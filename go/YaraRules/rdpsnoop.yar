/*
    This rule will look for common powershell elements
*/

rule rdp_snitch
{
    meta:
        author = "linted"
    strings:
        $json_start = /{/
        $ip = /"ip"/
        $asn = /"asn"/
        $isp = /isp/
        $org = /"org"/
        $regionName= /"regionName"/
        $country = /"country"/
        $account = /"account"/ 
        $keyboard = /"keyboard"/
        $client_build = /"client_build"/
        $client_name = /"client_name"/
        $ip_type = /"ip_type"/
    condition:
        all of them
}
