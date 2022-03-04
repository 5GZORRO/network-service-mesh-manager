package nsm

// States
var CREATING = "CREATING"

var WAIT_FOR_GATEWAY_CONFIG = "WAIT_FOR_GATEWAY_CONFIGURATION"
var CREATED = "CREATED"
var RUNNING = "RUNNING"
var READY = "READY"

// Intermediate states
var DELETING_RESOURCES = "DELETING_RESOURCES"
var CONFIGURING = "CONFIGURING"
var DELETING_CONFIGURATION = "DELETING_CONFIGURATION"

// Error states
var CREATION_ERROR = "CREATION_ERROR"
var CONFIGURATION_ERROR = "CONFIGURATION_ERROR"
var ERROR = "ERROR"
var DELETE_ERROR = "DELETE_ERROR"
