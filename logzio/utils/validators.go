package utils

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/logzio/logzio_terraform_client/alerts_v2"
	"github.com/logzio/logzio_terraform_client/archive_logs"
	"github.com/logzio/logzio_terraform_client/authentication_groups"
	"github.com/logzio/logzio_terraform_client/endpoints"
	"github.com/logzio/logzio_terraform_client/grafana_alerts"
	"github.com/logzio/logzio_terraform_client/s3_fetcher"
	"github.com/logzio/logzio_terraform_client/users"
	"regexp"
)

func contains(slice []string, s string) bool {

	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}

func ValidateOperationV2(v interface{}, k string) (ws []string, errors []error) {

	value := v.(string)

	validOperations := []string{
		alerts_v2.OperatorGreaterThanOrEquals,
		alerts_v2.OperatorGreaterThan,
		alerts_v2.OperatorEquals,
		alerts_v2.OperatorLessThan,
		alerts_v2.OperatorLessThanOrEquals,
		alerts_v2.OperatorNotEquals,
	}

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("operation %q must not be blank and be one of %s", k, validOperations))
	}

	valid := false
	for _, op := range validOperations {
		if op == value {
			valid = true
		}
	}

	if !valid {
		errors = append(errors, fmt.Errorf("operation %q must be one of %s", k, validOperations))
	}
	return
}

func ValidateOutputType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	validOutputTypes := []string{
		alerts_v2.OutputTypeJson,
		alerts_v2.OutputTypeTable,
	}

	if !contains(validOutputTypes, value) {
		errors = append(errors, fmt.Errorf("output type %q must be one of %s", k, validOutputTypes))
	}
	return
}

func ValidateSortTypes(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) == 0 {
		return
	}

	validTypes := []string{
		alerts_v2.SortAsc,
		alerts_v2.SortDesc,
	}

	if !contains(validTypes, value) {
		errors = append(errors, fmt.Errorf("sort type %q must be one of %s", k, validTypes))
	}

	return
}

func ValidateUrl(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	b, err := regexp.Match(VALIDATE_URL_REGEX, []byte(value))

	if !b && err == nil {
		err = fmt.Errorf("Bad URL provided")
	}

	if !b || err != nil {
		errors = append(errors, err)
	}

	return
}

func ValidateHttpMethod(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !findStringInArray(value, []string{"GET", "POST", "PUT", "DELETE"}) {
		errors = append(errors, fmt.Errorf("invalid HTTP method specified"))
	}

	return
}

func ValidateEndpointType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validTypes := []string{
		string(endpoints.EndpointTypeSlack),
		string(endpoints.EndpointTypeCustom),
		string(endpoints.EndpointTypePagerDuty),
		string(endpoints.EndpointTypeBigPanda),
		string(endpoints.EndpointTypeDataDog),
		string(endpoints.EndpointTypeVictorOps),
		string(endpoints.EndpointTypeServiceNow),
		string(endpoints.EndpointTypeOpsGenie),
		string(endpoints.EndpointTypeMicrosoftTeams)}

	if !contains(validTypes, value) {
		errors = append(errors, fmt.Errorf("value for endpoint type is unknown"))
	}

	return
}

func ValidateArchiveLogsStorageType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validStorageTypes := []string{
		archive_logs.StorageTypeS3,
		archive_logs.StorageTypeBlob,
	}

	if !contains(validStorageTypes, value) {
		errors = append(errors, fmt.Errorf("value for storage type is unknown. valid types are: %s", validStorageTypes))
	}

	return
}

func ValidateArchiveLogsAwsCredentialsType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validStorageTypes := []string{
		archive_logs.CredentialsTypeKeys,
		archive_logs.CredentialsTypeIam,
	}

	if !contains(validStorageTypes, value) {
		errors = append(errors, fmt.Errorf("value for credentials type is unknown. valid types are: %s", validStorageTypes))
	}

	return
}

func ValidateUserRole(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validUserRoles := []string{
		authentication_groups.AuthGroupsUserRoleRegular,
		authentication_groups.AuthGroupsUserRoleReadonly,
		authentication_groups.AuthGroupsUserRoleAdmin,
	}

	if !contains(validUserRoles, value) {
		errors = append(errors, fmt.Errorf("value for user role is unknown"))
	}

	return
}

func ValidateGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("group name must be set"))
	}

	return
}

func ValidateUserRoleUser(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	validUserRole := []string{
		users.UserRoleAccountAdmin,
		users.UserRoleRegular,
		users.UserRoleReadOnly,
	}

	if !contains(validUserRole, value) {
		errors = append(errors, fmt.Errorf("user role %q must be one of %s", k, validUserRole))
	}
	return
}

func ValidateScheduleTimezone(v interface{}, path cty.Path) diag.Diagnostics {
	timezone := v.(string)
	timezones := GetAlertV2ScheduleTimezones()
	if !contains(timezones, timezone) {
		return diag.Errorf("Timezone %s is not in the allowed timezones list.", timezone)
	}

	var diags diag.Diagnostics
	return diags
}

func ValidateS3FetcherRegion(v interface{}, path cty.Path) diag.Diagnostics {
	region := v.(string)
	regions := s3_fetcher.GetValidRegions()
	for _, validRegion := range regions {
		if region == validRegion.String() {
			return diag.Diagnostics{}
		}
	}

	return diag.Errorf("Region %s is not in the allowed aws regions list: %s", region, regions)
}

func ValidateS3FetcherLogsType(v interface{}, path cty.Path) diag.Diagnostics {
	logsType := v.(string)
	validLogsTypes := s3_fetcher.GetValidLogsType()
	for _, validType := range validLogsTypes {
		if logsType == validType.String() {
			return diag.Diagnostics{}
		}
	}

	return diag.Errorf("Logs type %s is not in the allowed logs types list: %s", logsType, validLogsTypes)
}

func ValidateExecErrState(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validExecErrState := []string{
		string(grafana_alerts.ErrAlerting),
		string(grafana_alerts.ErrOK),
		string(grafana_alerts.ErrError),
	}

	if !contains(validExecErrState, value) {
		errors = append(errors, fmt.Errorf("value for exec err state is unknown"))
	}

	return
}

func ValidateExecNoDataState(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validNoDataState := []string{
		string(grafana_alerts.NoDataAlerting),
		string(grafana_alerts.NoDataOk),
		string(grafana_alerts.NoData),
	}

	if !contains(validNoDataState, value) {
		errors = append(errors, fmt.Errorf("value for no data state is unknown"))
	}

	return
}
