/*
* This rule will look for syntax associated with configuration files routinely seen in Ubiquiti EdgeRouter Series (and likely others) 
* Tested against default configuration files found at https://github.com/stevejenkins/UBNT-EdgeRouter-Example-Configs
* as well as other configuration files posted on https://www.reddit.com/r/Ubiquiti
*/

rule er_configs
{
    meta:
        author = "calhandoh"
    strings:
        $firewall = /firewall {[\sa-zA-Z0-9-_{"\/.:$}]*/
        $interfaces = /interfaces {[\sa-zA-Z0-9-_{"\/.:$}]*/
        $system = /system {[\sa-zA-Z0-9-_{"\/.:$}]*/
    condition:
        // probably overkill as I only observed ($firewall, $interfaces) or ($firewall, $interfaces, $system) configs during testing
        all of ($firewall, $interfaces) or all of ($interfaces, $system) or all of ($firewall, $system) or all of ($firewall, $interfaces, $system)
}