package api

import (
	"database/sql"
	"net/http"

	"github.com/Juniper/contrail/pkg/common"
	"github.com/Juniper/contrail/pkg/generated/db"
	"github.com/Juniper/contrail/pkg/generated/models"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"
)

//KubernetesNodeRESTAPI
type KubernetesNodeRESTAPI struct {
	DB *sql.DB
}

type KubernetesNodeCreateRequest struct {
	Data *models.KubernetesNode `json:"kubernetes-node"`
}

//Path returns api path for collections.
func (api *KubernetesNodeRESTAPI) Path() string {
	return "/kubernetes-nodes"
}

//LongPath returns api path for elements.
func (api *KubernetesNodeRESTAPI) LongPath() string {
	return "/kubernetes-node/:id"
}

//SetDB sets db object
func (api *KubernetesNodeRESTAPI) SetDB(db *sql.DB) {
	api.DB = db
}

//Create handle a Create REST API.
func (api *KubernetesNodeRESTAPI) Create(c echo.Context) error {
	requestData := &KubernetesNodeCreateRequest{
		Data: models.MakeKubernetesNode(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "kubernetes_node",
		}).Debug("bind failed on create")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	model := requestData.Data
	if model == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	if model.UUID == "" {
		model.UUID = uuid.NewV4().String()
	}
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.CreateKubernetesNode(tx, model)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "kubernetes_node",
		}).Debug("db create failed on create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusCreated, requestData)
}

//Update handles a REST Update request.
func (api *KubernetesNodeRESTAPI) Update(c echo.Context) error {
	return nil
}

//Delete handles a REST Delete request.
func (api *KubernetesNodeRESTAPI) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.DeleteKubernetesNode(tx, id)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//Show handles a REST Show request.
func (api *KubernetesNodeRESTAPI) Show(c echo.Context) error {
	id := c.Param("id")
	var result *models.KubernetesNode
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ShowKubernetesNode(tx, id)
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"kubernetes_node": result,
	})
}

//List handles a List REST API Request.
func (api *KubernetesNodeRESTAPI) List(c echo.Context) error {
	var result []*models.KubernetesNode
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListKubernetesNode(tx, &common.ListSpec{
				Limit: 1000,
			})
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"kubernetes-nodes": result,
	})
}
