#compdef _prog prog


function _prog {
  local -a commands

  _arguments -C \
    '--abool[A bool flag]' \
    '--abytes[A bytes base64 flag]:' \
    '--aduration[A duration flag]:' \
    '--afloat[A float flag. THIS IS A VIPER BUG! Gets transformed to a string!]:' \
    '--anint[An int flag]:' \
    '--anip[An IP flag]:' \
    '--anipnet[An IPNet flag]:' \
    '*--astringslice[A string slice flag]:' \
    '(-c --config)'{-c,--config}'[Load configuration from this file. The path may be absolute or relative. Supported extensions: json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env, ini]:' \
    '--debug-address[Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.]:' \
    '--debug-autostart[Whether the golang pprof server is automatically started.]' \
    '--log-file-path[Path to log file. Only required if --log-target is set to file.]:' \
    '--log-format[Define default log format, one of: json, keyvalue, pretty.]:' \
    '--log-level[Define minimum log level, one of: panic, fatal, error, warning, info, debug, trace.]:' \
    '--log-target[Define log output target, one of: stdout, stderr, file.]:' \
    '--remote-first[A flag which is in in a subtree]:' \
    '--remote-second[The other part of the subtree flag so we can see what this means]:' \
    '(-s --short)'{-s,--short}'[A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?]:' \
    '--store-another-sub-level[A flag with a second sublevel, sometimes with dashes]:' \
    '(-S --store-dir)'{-S,--store-dir}'[A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree]:' \
    '(-L --store-log-level)'{-L,--store-log-level}'[The other subtree element has a dash as well]:' \
    '--version[Print version information]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "completion:Generate a completion script"
      "config:Manage configuration options"
      "help:Help about any command"
      "version:Print version information"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  completion)
    _prog_completion
    ;;
  config)
    _prog_config
    ;;
  help)
    _prog_help
    ;;
  version)
    _prog_version
    ;;
  esac
}


function _prog_completion {
  local -a commands

  _arguments -C \
    '--abool[A bool flag]' \
    '--abytes[A bytes base64 flag]:' \
    '--aduration[A duration flag]:' \
    '--afloat[A float flag. THIS IS A VIPER BUG! Gets transformed to a string!]:' \
    '--anint[An int flag]:' \
    '--anip[An IP flag]:' \
    '--anipnet[An IPNet flag]:' \
    '*--astringslice[A string slice flag]:' \
    '(-c --config)'{-c,--config}'[Load configuration from this file. The path may be absolute or relative. Supported extensions: json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env, ini]:' \
    '--debug-address[Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.]:' \
    '--debug-autostart[Whether the golang pprof server is automatically started.]' \
    '--log-file-path[Path to log file. Only required if --log-target is set to file.]:' \
    '--log-format[Define default log format, one of: json, keyvalue, pretty.]:' \
    '--log-level[Define minimum log level, one of: panic, fatal, error, warning, info, debug, trace.]:' \
    '--log-target[Define log output target, one of: stdout, stderr, file.]:' \
    '--remote-first[A flag which is in in a subtree]:' \
    '--remote-second[The other part of the subtree flag so we can see what this means]:' \
    '(-s --short)'{-s,--short}'[A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?]:' \
    '--store-another-sub-level[A flag with a second sublevel, sometimes with dashes]:' \
    '(-S --store-dir)'{-S,--store-dir}'[A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree]:' \
    '(-L --store-log-level)'{-L,--store-log-level}'[The other subtree element has a dash as well]:' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "bash:Generates bash completion script"
      "powershell:Generates powershell completion script"
      "zsh:Generates zsh completion script"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  bash)
    _prog_completion_bash
    ;;
  powershell)
    _prog_completion_powershell
    ;;
  zsh)
    _prog_completion_zsh
    ;;
  esac
}

function _prog_completion_bash {
  _arguments \
    '(-h --help)'{-h,--help}'[help for bash]' \
    '--abool[A bool flag]' \
    '--abytes[A bytes base64 flag]:' \
    '--aduration[A duration flag]:' \
    '--afloat[A float flag. THIS IS A VIPER BUG! Gets transformed to a string!]:' \
    '--anint[An int flag]:' \
    '--anip[An IP flag]:' \
    '--anipnet[An IPNet flag]:' \
    '*--astringslice[A string slice flag]:' \
    '(-c --config)'{-c,--config}'[Load configuration from this file. The path may be absolute or relative. Supported extensions: json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env, ini]:' \
    '--debug-address[Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.]:' \
    '--debug-autostart[Whether the golang pprof server is automatically started.]' \
    '--log-file-path[Path to log file. Only required if --log-target is set to file.]:' \
    '--log-format[Define default log format, one of: json, keyvalue, pretty.]:' \
    '--log-level[Define minimum log level, one of: panic, fatal, error, warning, info, debug, trace.]:' \
    '--log-target[Define log output target, one of: stdout, stderr, file.]:' \
    '--remote-first[A flag which is in in a subtree]:' \
    '--remote-second[The other part of the subtree flag so we can see what this means]:' \
    '(-s --short)'{-s,--short}'[A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?]:' \
    '--store-another-sub-level[A flag with a second sublevel, sometimes with dashes]:' \
    '(-S --store-dir)'{-S,--store-dir}'[A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree]:' \
    '(-L --store-log-level)'{-L,--store-log-level}'[The other subtree element has a dash as well]:'
}

function _prog_completion_powershell {
  _arguments \
    '--abool[A bool flag]' \
    '--abytes[A bytes base64 flag]:' \
    '--aduration[A duration flag]:' \
    '--afloat[A float flag. THIS IS A VIPER BUG! Gets transformed to a string!]:' \
    '--anint[An int flag]:' \
    '--anip[An IP flag]:' \
    '--anipnet[An IPNet flag]:' \
    '*--astringslice[A string slice flag]:' \
    '(-c --config)'{-c,--config}'[Load configuration from this file. The path may be absolute or relative. Supported extensions: json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env, ini]:' \
    '--debug-address[Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.]:' \
    '--debug-autostart[Whether the golang pprof server is automatically started.]' \
    '--log-file-path[Path to log file. Only required if --log-target is set to file.]:' \
    '--log-format[Define default log format, one of: json, keyvalue, pretty.]:' \
    '--log-level[Define minimum log level, one of: panic, fatal, error, warning, info, debug, trace.]:' \
    '--log-target[Define log output target, one of: stdout, stderr, file.]:' \
    '--remote-first[A flag which is in in a subtree]:' \
    '--remote-second[The other part of the subtree flag so we can see what this means]:' \
    '(-s --short)'{-s,--short}'[A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?]:' \
    '--store-another-sub-level[A flag with a second sublevel, sometimes with dashes]:' \
    '(-S --store-dir)'{-S,--store-dir}'[A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree]:' \
    '(-L --store-log-level)'{-L,--store-log-level}'[The other subtree element has a dash as well]:'
}

function _prog_completion_zsh {
  _arguments \
    '(-h --help)'{-h,--help}'[help for zsh]' \
    '--abool[A bool flag]' \
    '--abytes[A bytes base64 flag]:' \
    '--aduration[A duration flag]:' \
    '--afloat[A float flag. THIS IS A VIPER BUG! Gets transformed to a string!]:' \
    '--anint[An int flag]:' \
    '--anip[An IP flag]:' \
    '--anipnet[An IPNet flag]:' \
    '*--astringslice[A string slice flag]:' \
    '(-c --config)'{-c,--config}'[Load configuration from this file. The path may be absolute or relative. Supported extensions: json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env, ini]:' \
    '--debug-address[Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.]:' \
    '--debug-autostart[Whether the golang pprof server is automatically started.]' \
    '--log-file-path[Path to log file. Only required if --log-target is set to file.]:' \
    '--log-format[Define default log format, one of: json, keyvalue, pretty.]:' \
    '--log-level[Define minimum log level, one of: panic, fatal, error, warning, info, debug, trace.]:' \
    '--log-target[Define log output target, one of: stdout, stderr, file.]:' \
    '--remote-first[A flag which is in in a subtree]:' \
    '--remote-second[The other part of the subtree flag so we can see what this means]:' \
    '(-s --short)'{-s,--short}'[A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?]:' \
    '--store-another-sub-level[A flag with a second sublevel, sometimes with dashes]:' \
    '(-S --store-dir)'{-S,--store-dir}'[A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree]:' \
    '(-L --store-log-level)'{-L,--store-log-level}'[The other subtree element has a dash as well]:'
}


function _prog_config {
  local -a commands

  _arguments -C \
    '--abool[A bool flag]' \
    '--abytes[A bytes base64 flag]:' \
    '--aduration[A duration flag]:' \
    '--afloat[A float flag. THIS IS A VIPER BUG! Gets transformed to a string!]:' \
    '--anint[An int flag]:' \
    '--anip[An IP flag]:' \
    '--anipnet[An IPNet flag]:' \
    '*--astringslice[A string slice flag]:' \
    '(-c --config)'{-c,--config}'[Load configuration from this file. The path may be absolute or relative. Supported extensions: json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env, ini]:' \
    '--debug-address[Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.]:' \
    '--debug-autostart[Whether the golang pprof server is automatically started.]' \
    '--log-file-path[Path to log file. Only required if --log-target is set to file.]:' \
    '--log-format[Define default log format, one of: json, keyvalue, pretty.]:' \
    '--log-level[Define minimum log level, one of: panic, fatal, error, warning, info, debug, trace.]:' \
    '--log-target[Define log output target, one of: stdout, stderr, file.]:' \
    '--remote-first[A flag which is in in a subtree]:' \
    '--remote-second[The other part of the subtree flag so we can see what this means]:' \
    '(-s --short)'{-s,--short}'[A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?]:' \
    '--store-another-sub-level[A flag with a second sublevel, sometimes with dashes]:' \
    '(-S --store-dir)'{-S,--store-dir}'[A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree]:' \
    '(-L --store-log-level)'{-L,--store-log-level}'[The other subtree element has a dash as well]:' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "show:Display the currently loaded configuration in specified format"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  show)
    _prog_config_show
    ;;
  esac
}

function _prog_config_show {
  _arguments \
    '--abool[A bool flag]' \
    '--abytes[A bytes base64 flag]:' \
    '--aduration[A duration flag]:' \
    '--afloat[A float flag. THIS IS A VIPER BUG! Gets transformed to a string!]:' \
    '--anint[An int flag]:' \
    '--anip[An IP flag]:' \
    '--anipnet[An IPNet flag]:' \
    '*--astringslice[A string slice flag]:' \
    '(-c --config)'{-c,--config}'[Load configuration from this file. The path may be absolute or relative. Supported extensions: json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env, ini]:' \
    '--debug-address[Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.]:' \
    '--debug-autostart[Whether the golang pprof server is automatically started.]' \
    '--log-file-path[Path to log file. Only required if --log-target is set to file.]:' \
    '--log-format[Define default log format, one of: json, keyvalue, pretty.]:' \
    '--log-level[Define minimum log level, one of: panic, fatal, error, warning, info, debug, trace.]:' \
    '--log-target[Define log output target, one of: stdout, stderr, file.]:' \
    '--remote-first[A flag which is in in a subtree]:' \
    '--remote-second[The other part of the subtree flag so we can see what this means]:' \
    '(-s --short)'{-s,--short}'[A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?]:' \
    '--store-another-sub-level[A flag with a second sublevel, sometimes with dashes]:' \
    '(-S --store-dir)'{-S,--store-dir}'[A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree]:' \
    '(-L --store-log-level)'{-L,--store-log-level}'[The other subtree element has a dash as well]:' \
    '1: :("dotenv" "env" "hcl" "ini" "json" "prop" "properties" "props" "toml" "yaml" "yml")'
}

function _prog_help {
  _arguments \
    '--abool[A bool flag]' \
    '--abytes[A bytes base64 flag]:' \
    '--aduration[A duration flag]:' \
    '--afloat[A float flag. THIS IS A VIPER BUG! Gets transformed to a string!]:' \
    '--anint[An int flag]:' \
    '--anip[An IP flag]:' \
    '--anipnet[An IPNet flag]:' \
    '*--astringslice[A string slice flag]:' \
    '(-c --config)'{-c,--config}'[Load configuration from this file. The path may be absolute or relative. Supported extensions: json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env, ini]:' \
    '--debug-address[Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.]:' \
    '--debug-autostart[Whether the golang pprof server is automatically started.]' \
    '--log-file-path[Path to log file. Only required if --log-target is set to file.]:' \
    '--log-format[Define default log format, one of: json, keyvalue, pretty.]:' \
    '--log-level[Define minimum log level, one of: panic, fatal, error, warning, info, debug, trace.]:' \
    '--log-target[Define log output target, one of: stdout, stderr, file.]:' \
    '--remote-first[A flag which is in in a subtree]:' \
    '--remote-second[The other part of the subtree flag so we can see what this means]:' \
    '(-s --short)'{-s,--short}'[A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?]:' \
    '--store-another-sub-level[A flag with a second sublevel, sometimes with dashes]:' \
    '(-S --store-dir)'{-S,--store-dir}'[A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree]:' \
    '(-L --store-log-level)'{-L,--store-log-level}'[The other subtree element has a dash as well]:'
}

function _prog_version {
  _arguments \
    '--abool[A bool flag]' \
    '--abytes[A bytes base64 flag]:' \
    '--aduration[A duration flag]:' \
    '--afloat[A float flag. THIS IS A VIPER BUG! Gets transformed to a string!]:' \
    '--anint[An int flag]:' \
    '--anip[An IP flag]:' \
    '--anipnet[An IPNet flag]:' \
    '*--astringslice[A string slice flag]:' \
    '(-c --config)'{-c,--config}'[Load configuration from this file. The path may be absolute or relative. Supported extensions: json, toml, yaml, yml, properties, props, prop, hcl, dotenv, env, ini]:' \
    '--debug-address[Listen address to use for the golang pprof server. Send USR2 signal to the process to toggle the server on and off.]:' \
    '--debug-autostart[Whether the golang pprof server is automatically started.]' \
    '--log-file-path[Path to log file. Only required if --log-target is set to file.]:' \
    '--log-format[Define default log format, one of: json, keyvalue, pretty.]:' \
    '--log-level[Define minimum log level, one of: panic, fatal, error, warning, info, debug, trace.]:' \
    '--log-target[Define log output target, one of: stdout, stderr, file.]:' \
    '--remote-first[A flag which is in in a subtree]:' \
    '--remote-second[The other part of the subtree flag so we can see what this means]:' \
    '(-s --short)'{-s,--short}'[A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?]:' \
    '--store-another-sub-level[A flag with a second sublevel, sometimes with dashes]:' \
    '(-S --store-dir)'{-S,--store-dir}'[A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree]:' \
    '(-L --store-log-level)'{-L,--store-log-level}'[The other subtree element has a dash as well]:'
}

