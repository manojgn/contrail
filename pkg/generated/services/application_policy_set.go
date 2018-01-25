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

//RESTCreateApplicationPolicySet handle a Create REST service.
func (service *ContrailService) RESTCreateApplicationPolicySet(c echo.Context) error {
    requestData := &models.ApplicationPolicySetCreateRequest{
        ApplicationPolicySet: models.MakeApplicationPolicySet(),
    }
    if err := c.Bind(requestData); err != nil {
        log.WithFields(log.Fields{
            "err": err,
            "resource": "application_policy_set",
        }).Debug("bind failed on create")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
    ctx := c.Request().Context()
    response, err := service.CreateApplicationPolicySet(ctx, requestData)
    if err != nil {
        return common.ToHTTPError(err)
    } 
    return c.JSON(http.StatusCreated, response)
}

//CreateApplicationPolicySet handle a Create API
func (service *ContrailService) CreateApplicationPolicySet(
    ctx context.Context, 
    request *models.ApplicationPolicySetCreateRequest) (*models.ApplicationPolicySetCreateResponse, error) {
    model := request.ApplicationPolicySet
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
            return db.CreateApplicationPolicySet(tx, model)
        }); err != nil {
        log.WithFields(log.Fields{
            "err": err,
            "resource": "application_policy_set",
        }).Debug("db create failed on create")
       return nil, common.ErrorInternal 
    }
    return &models.ApplicationPolicySetCreateResponse{
        ApplicationPolicySet: request.ApplicationPolicySet,
    }, nil
}

//RESTUpdateApplicationPolicySet handles a REST Update request.
func (service *ContrailService) RESTUpdateApplicationPolicySet(c echo.Context) error {
    id := c.Param("id")
    request := &models.ApplicationPolicySetUpdateRequest{}
    if err := c.Bind(request); err != nil {
            log.WithFields(log.Fields{
                "err": err,
                "resource": "application_policy_set",
            }).Debug("bind failed on update")
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
    }
    request.ID = id
    ctx := c.Request().Context()
    response, err := service.UpdateApplicationPolicySet(ctx, request)
    if err != nil {
        return nil, common.ToHTTPError(err)
    }
    return c.JSON(http.StatusOK, response)
}

//UpdateApplicationPolicySet handles a Update request.
func (service *ContrailService) UpdateApplicationPolicySet(ctx context.Context, request *models.ApplicationPolicySetUpdateRequest) (*models.ApplicationPolicySetUpdateResponse, error) {
    id = request.ID
    model = request.ApplicationPolicySet
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
            return db.UpdateApplicationPolicySet(tx, id, model)
        }); err != nil {
        log.WithFields(log.Fields{
            "err": err,
            "resource": "application_policy_set",
        }).Debug("db update failed")
        return nil, common.ErrorInternal
    }
    return &models.ApplicationPolicySet.UpdateResponse{
        ApplicationPolicySet: model,
    }, nil
}

//RESTDeleteApplicationPolicySet delete a resource using REST service.
func (service *ContrailService) RESTDeleteApplicationPolicySet(c echo.Context) error {
    id := c.Param("id")
    request := &models.ApplicationPolicySetDeleteRequest{
        ID: id
    } 
    ctx := c.Request().Context()
    response, err := service.DeleteApplicationPolicySet(ctx, request)
    if err != nil {
        return common.ToHTTPError(err)
    }
    return c.JSON(http.StatusNoContent, nil)
}

//DeleteApplicationPolicySet delete a resource.
func (service *ContrailService) DeleteApplicationPolicySet(ctx context.Context, request *models.ApplicationPolicySetDeleteRequest) (*models.ApplicationPolicySetDeleteResponse, error) {
    id := request.ID
    auth := common.GetAuthCTX(ctx)
    if err := common.DoInTransaction(
        service.DB,
        func (tx *sql.Tx) error {
            return db.DeleteApplicationPolicySet(tx, id, auth)
        }); err != nil {
            log.WithField("err", err).Debug("error deleting a resource")
        return nil, common.ErrorInternal
    }
    return &models.ApplicationPolicySetDeleteResponse{
        ID: id,
    }, nil
}

//RESTShowApplicationPolicySet a REST Show request.
func (service *ContrailService) RESTShowApplicationPolicySet(c echo.Context) (error) {
    id := c.Param("id")
    auth := common.GetAuthContext(c)
    var result []*models.ApplicationPolicySet
    var err error
    if err := common.DoInTransaction(
        service.DB,
        func (tx *sql.Tx) error {
            result, err = db.ListApplicationPolicySet(tx, &common.ListSpec{
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
        "application_policy_set": result,
    })
}

//RESTListApplicationPolicySet handles a List REST service Request.
func (service *ContrailService) RESTListApplicationPolicySet(c echo.Context) (error) {
    var result []*models.ApplicationPolicySet
    var err error
    auth := common.GetAuthContext(c)
    listSpec := common.GetListSpec(c)
    listSpec.Auth = auth
    if err := common.DoInTransaction(
        service.DB,
        func (tx *sql.Tx) error {
            result, err = db.ListApplicationPolicySet(tx, listSpec)
            return err
        }); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
    }
    return c.JSON(http.StatusOK, map[string]interface{}{
        "application-policy-sets": result,
    })
}