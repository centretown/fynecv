package svc

import "fmt"

var (
	serviceFormat = `{ "type":"call_service", "domain":"%s", "service":"%s", "service_data": {%s}, "target":{"entity_id":"%s"},`
	idstr         = `"id":%d }`
)

func ServiceCmd(domain string, service string, entityID string, data map[string]string) string {
	var serviceData string
	for k, v := range data {
		serviceData += fmt.Sprintf(`"%s":%s,`, k, v)
	}
	if len(serviceData) > 0 {
		serviceData = serviceData[:len(serviceData)-1]
	}
	cmd := fmt.Sprintf(serviceFormat, domain, service, serviceData, entityID) + idstr
	return cmd
}

//cmd := svc.ServiceCmd("light", "turn_on", entityID, d)

func LightCmd(entityID string, data map[string]string) string {
	return ServiceCmd("light", "turn_on", entityID, data)
}
