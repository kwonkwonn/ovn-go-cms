// Code generated by "libovsdb.modelgen"
// DO NOT EDIT.

package NBModel

const GatewayChassisTable = "Gateway_Chassis"

// GatewayChassis defines an object in Gateway_Chassis table
type GatewayChassis struct {
	UUID        string            `ovsdb:"_uuid"`
	ChassisName string            `ovsdb:"chassis_name"`
	ExternalIDs map[string]string `ovsdb:"external_ids"`
	Name        string            `ovsdb:"name"`
	Options     map[string]string `ovsdb:"options"`
	Priority    int               `ovsdb:"priority"`
}
