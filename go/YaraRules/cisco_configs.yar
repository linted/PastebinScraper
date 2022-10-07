/*
    This rule will look for Cisco configuration files
*/

rule cisco_configs
{
    meta:
        author = "@calhandoh"
    strings:
        $config_start = "Building configuration..." wide ascii nocase
        $version = "version" wide ascii nocase
        $interface = "interface" wide ascii nocase
        $current_config = "Current configuration" wide ascii nocase
        $line_break = /!\s/ wide ascii nocase
    condition:
        all of them
}