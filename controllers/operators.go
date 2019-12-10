package controllers

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/avast/retry-go"
	"github.com/go-logr/logr"
	apiv1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OperatorMode interface for operating mode (controller|sidecar)
type OperatorMode interface {
	CreateOrUpdate(ctx context.Context, oathkeeperRulesJSON []byte, nameFunc func() types.NamespacedName)
	GetConfigLocation() types.NamespacedName
}

type ConfigMapOperator struct {
	client.Client
	Log           logr.Logger
	RuleConfigmap types.NamespacedName
	RulesFileName string
}

type FilesOperator struct {
	Log           logr.Logger
	RulesFilePath string
}

func (cmo *ConfigMapOperator) updateOrCreateRulesConfigmap(ctx context.Context, configMap types.NamespacedName, data string) error {

	var oathkeeperRulesConfigmap apiv1.ConfigMap
	var exists = false

	fetchMapFunc := func() error {

		if err := cmo.Get(ctx, configMap, &oathkeeperRulesConfigmap); err != nil {

			if apierrs.IsForbidden(err) {
				return retry.Unrecoverable(err)
			}

			if apierrs.IsNotFound(err) {
				return nil
			}

			return err
		}

		exists = true
		return nil
	}

	createMapFunc := func() error {
		cmo.Log.Info("creating ConfigMap")
		oathkeeperRulesConfigmap = apiv1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMap.Name,
				Namespace: configMap.Namespace,
			},
			Data: map[string]string{cmo.RulesFileName: data},
		}
		return cmo.Create(ctx, &oathkeeperRulesConfigmap)
	}

	updateMapFunc := func() error {
		cmo.Log.Info("updating ConfigMap")
		oathkeeperRulesConfigmap.Data = map[string]string{cmo.RulesFileName: data}
		err := cmo.Update(ctx, &oathkeeperRulesConfigmap)
		return err
	}

	return retryOnError(func() error {
		exists = false

		if err := fetchMapFunc(); err != nil {
			return err
		}

		if exists {
			err := updateMapFunc()
			if err != nil {
				if isObjectHasBeenModified(err) {
					cmo.Log.Error(err, "incorrect object version during ConfigMap update")
				}
			}
			return err
		}

		return createMapFunc()
	})
}

func (cmo *ConfigMapOperator) GetConfigLocation() types.NamespacedName {
	return cmo.RuleConfigmap
}

func (cmo *ConfigMapOperator) CreateOrUpdate(ctx context.Context, oathkeeperRulesJSON []byte, nameFunc func() types.NamespacedName) {
	configMap := nameFunc()
	if err := cmo.updateOrCreateRulesConfigmap(ctx, configMap, string(oathkeeperRulesJSON)); err != nil {
		cmo.Log.Error(err, "unable to process rules Configmap")
		os.Exit(1)
	}
}

func (fo *FilesOperator) updateOrCreateRulesFile(ctx context.Context, data string) error {
	var f *os.File
	f, err := os.Create(fo.RulesFilePath)
	if err != nil {
		fo.Log.Error(err, "error while creating config file")
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	byteCount, err := w.WriteString(data)
	fo.Log.Info(fmt.Sprintf("wiriting %d bytes of data into %s", byteCount, fo.RulesFilePath))
	w.Flush()
	if err != nil {
		fo.Log.Error(err, "error while writing to file")
		return err
	}
	return nil
}

func (fo *FilesOperator) CreateOrUpdate(ctx context.Context, oathkeeperRulesJSON []byte, nameFunc func() types.NamespacedName) {
	err := fo.updateOrCreateRulesFile(ctx, string(oathkeeperRulesJSON))
	if err != nil {
		fo.Log.Error(err, "unable to process rules Configmap")
		os.Exit(1)
	}
}

func (cmo *FilesOperator) GetConfigLocation() types.NamespacedName {
	return types.NamespacedName{
		Name:      cmo.RulesFilePath,
		Namespace: "local_file",
	}
}
