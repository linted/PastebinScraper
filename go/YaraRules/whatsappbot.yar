/*
    Whats app bot creds
*/

rule whatsappcreds
{
    meta:
        author = "linted"
        info = "Whats app bot creds"

    strings:
        $creds = /{"creds":{"noiseKey":{"private":{"type":"Buffer","data":/
    condition:
        $creds

}
