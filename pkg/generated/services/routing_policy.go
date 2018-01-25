package services 

import (
    "context"
    "net/http"
    "database/sql"
    "github.com/Juniper/contrail/pkg/generated/models"
    "github.com/Juniper/contrail/pkg/generated/db"
    "github.com/satori/go.uuid"
    "github.com/labstack/echo"
    "github.com/Juniper/contrail/pkg/common"

	log "github.com/sirupsen/logrus"
)

//RESTCreateRoutingPolicy handle a Create REST service.
func (service *ContrailService) RESTCreateRoutingPolicy(c echo.Context) error {
    requestData := &models.RoutingPolicyCreateRequest{
        RoutingPolicy: models.MakeRoutingPolicy(),
    }
    if err := c.Bind(requestData); err != nil {
        log.WithFields(log.Fields{
            "err": err,
            "resource": "routing_policy",
        }).Debug("bind failed on create")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
    ctx := c.Request().Context()
    response, err := service.CreateRoutingPolicy(ctx, requestData)
    if err != nil {
        return common.ToHTTPError(err)
    } 
    return c.JSON(http.StatusCreated, response)
}

//CreateRoutingPolicy handle a Create API
func (service *ContrailService) CreateRoutingPolicy(
    ctx context.Context, 
    request *models.RoutingPolicyCreateRequest) (*models.RoutingPolicyCreateResponse, error) {
    model := request.RoutingPolicy
    if model.UUID == "" {
        model.UUID = uuid.NewV4().String()
    }

    if model.FQName == nil {
       return nil, common.ErrorBadRequest("Missing fq_name")
    }

    auth := common.GetAuthCTX(ctx)
    if auth == nil {
        return nil, common.ErrorUnauthenticated
    }
    model.Perms2.Owner = auth.ProjectID()
    if err := common.DoInTransaction(
        service.DB,
        func (tx *sql.Tx) error {
            return db.CreateRoutingPolicy(tx, model)
        }); err != nil {
        log.WithFields(log.Fields{
            "err": err,
            "resource": "routing_policy",
        }).Debug("db create failed on create")
       return nil, common.ErrorInternal 
    }
    return &models.RoutingPolicyCreateResponse{
        RoutingPolicy: request.RoutingPolicy,
    }, nil
}

//RESTUpdateRoutingPolicy handles a REST Update request.
func (service *ContrailService) RESTUpdateRoutingPolicy(c echo.Context) error {
    id := c.Param("id")
    request := &models.RoutingPolicyUpdateRequest{}
    if err := c.Bind(request); err != nil {
            log.WithFields(log.Fields{
                "err": err,
                "resource": "routing_policy",
            }).Debug("bind failed on update")
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
    }
    request.ID = id
    ctx := c.Request().Context()
    response, err := service.UpdateRoutingPolicy(ctx, request)
    if err != nil {
        return nil, common.ToHTTPError(err)
    }
    return c.JSON(http.StatusOK, response)
}

//UpdateRoutingPolicy handles a Update request.
func (service *ContrailService) UpdateRoutingPolicy(ctx context.Context, request *models.RoutingPolicyUpdateRequest) (*models.RoutingPolicyUpdateResponse, error) {
    id = request.ID
    model = request.RoutingPolicy
    if model == nil {
        return nil, common.ErrorBadRequest("Update body is empty")
    }
    auth := common.GetAuthCTX(ctx)
    ok := common.SetValueByPath(model, "Perms2.Owner", ".", auth.ProjectID())
    if !ok {
        return nil, common.ErrorBadRequest("Invalid JSON format")
    }
    if err := common.DoInTransaction(
        service.DB,
        func (tx *sql.Tx) error {
            return db.UpdateRoutingPolicy(tx, id, model)
        }); err != nil {
        log.WithFields(log.Fields{
            "err": err,
            "resource": "routing_policy",
        }).Debug("db update failed")
        return nil, common.ErrorInternal
    }
    return &models.RoutingPolicy.UpdateResponse{
        RoutingPolicy: model,
    }, nil
}

//RESTDeleteRoutingPolicy delete a resource using REST service.
func (service *ContrailService) RESTDeleteRoutingPolicy(c echo.Context) error {
    id := c.Param("id")
    request := &models.RoutingPolicyDeleteRequest{
        ID: id
    } 
    ctx := c.Request().Context()
    response, err := service.DeleteRoutingPolicy(ctx, request)
    if err != nil {
        return common.ToHTTPError(err)
    }
    return c.JSON(http.StatusNoContent, nil)
}

//DeleteRoutingPolicy delete a resource.
func (service *ContrailService) DeleteRoutingPolicy(ctx context.Context, request *models.RoutingPolicyDeleteRequest) (*models.RoutingPolicyDeleteResponse, error) {
    id := request.ID
    auth := common.GetAuthCTX(ctx)
    if err := common.DoInTransaction(
        service.DB,
        func (tx *sql.Tx) error {
            return db.DeleteRoutingPolicy(tx, id, auth)
        }); err != nil {
            log.WithField("err", err).Debug("error deleting a resource")
        return nil, common.ErrorInternal
    }
    return &models.RoutingPolicyDeleteResponse{
        ID: id,
    }, nil
}

//RESTShowRoutingPolicy a REST Show request.
func (service *ContrailService) RESTShowRoutingPolicy(c echo.Context) (error) {
    id := c.Param("id")
    auth := common.GetAuthContext(c)
    var result []*models.RoutingPolicy
    var err error
    if err := common.DoInTransaction(
        service.DB,
        func (tx *sql.Tx) error {
            result, err = db.ListRoutingPolicy(tx, &common.ListSpec{
                Limit: 1,
                Auth: auth,
                Filter: common.Filter{
                    "uuid": []string{id},
                },
            })
            return err
        }); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
    }
    return c.JSON(http.StatusOK, map[string]interface{}{
        "routing_policy": result,
    })
}

//RESTListRoutingPolicy handles a List REST service Request.
func (service *ContrailService) RESTListRoutingPolicy(c echo.Context) (error) {
    var result []*models.RoutingPolicy
    var err error
    auth := common.GetAuthContext(c)
    listSpec := common.GetListSpec(c)
    listSpec.Auth = auth
    if err := common.DoInTransaction(
        service.DB,
        func (tx *sql.Tx) error {
            result, err = db.ListRoutingPolicy(tx, listSpec)
            return err
        }); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
    }
    return c.JSON(http.StatusOK, map[string]interface{}{
        "routing-policys": result,
    })
}