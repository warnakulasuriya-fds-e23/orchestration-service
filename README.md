# Environmental Variable Configuration
following environment variables are required for the orchestration service to function
```
ADRESS_FOR_OUTGOING (eg := https://localhost:9443)
CONSUMER_KEY_FOR_OUTGOING (eg := wCPmxlccE9w1lajdflWURD6fuTqWdBhjadjf3PIa)
CONSUMER_SECRET_FOR_OUTGOING (eg := xij6lkadftnP_n_QV_I8vadjoPtIytMoifHtA67E1h8qqRHtYr8Aa)
TOKEN_ENDPOINT_FOR_OUTGOING (eg := https://localhost:9443/oauth2/token)
IDP_APP_REDIRECT_URI (eg := http://localhost:4000)
IDP_BIO_SDK_AUTHENTICATOR_DISPLAY_NAME (eg := BioSDKService)
DEVICES_CONFIG_JSON_PATH (eg := ./devicesConfig.json)
```
ADRESS_FOR_OUTGOING is the adress of the WSO2 IDP being used

CONSUMER_KEY_FOR_OUTGOING and CONSUMER_SECRET_FOR_OUTGOING are for oauth access token generation for the calls that will be sent to the WSO2 IDP (an OAuth2.0 supporting application needs to be setup in the IDP)

TOKEN_ENDPOINT_FOR_OUTGOING is the token endpoint that is specified in the info tab of the configured OAuth2.0 supporting application of the WSO2 IDP being used

IDP_APP_REDIRECT_URI is the redirect url specified in the the configured OAuth2.0 supporting application of the WSO2 IDP being used

IDP_BIO_SDK_AUTHENTICATOR_DISPLAY_NAME is the display name of the custom authentication service setup using the connections tab of the relevant WSO2 IDP being used (note that this custom authentication should be added to the first execution stage of the login flow of the configured OAuth2.0 supporting application, the orchestration service will find the authenticator then as an available authenticator even if other authenticators are setup in first execution stage. )

DEVICES_CONFIG_JSON_PATH is the path of the json file containing biometric devices that should be registered with the orchestration service.

# Sample Devices Config json
follow is a sample for the json file that should be given the path to in DEVICES_CONFIG_JSON_PATH.
```
{
  "enrollmentDevices": [
    "1234ABCD"
  ],
  "deviceDetails": {
    "RNYRI378": {
      "floor": "First Floor",
      "door": "Main Entrance"
    },
    "ABCD123": {
      "floor": "First Floor",
      "door": "Secondary Entrance"
    },
    "XSA4242": {
        "floor": "Second Floor",
        "door" : "Main Entrance"
    }
  }
}
```

## Note
the configured OAuth2.0 supporting application of the WSO2 IDP being used should have a MFA login flow. The Bio SDK service utilization is designed be utilized in the first execution stage itself


