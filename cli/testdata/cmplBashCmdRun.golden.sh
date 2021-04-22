# bash completion for prog                                 -*- shell-script -*-

__prog_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__prog_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__prog_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__prog_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__prog_handle_go_custom_completion()
{
    __prog_debug "${FUNCNAME[0]}: cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}"

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling ${words[0]} instead of directly prog allows to handle aliases
    args=("${words[@]:1}")
    requestComp="${words[0]} __completeNoDesc ${args[*]}"

    lastParam=${words[$((${#words[@]}-1))]}
    lastChar=${lastParam:$((${#lastParam}-1)):1}
    __prog_debug "${FUNCNAME[0]}: lastParam ${lastParam}, lastChar ${lastChar}"

    if [ -z "${cur}" ] && [ "${lastChar}" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __prog_debug "${FUNCNAME[0]}: Adding extra empty parameter"
        requestComp="${requestComp} \"\""
    fi

    __prog_debug "${FUNCNAME[0]}: calling ${requestComp}"
    # Use eval to handle any environment variables and such
    out=$(eval "${requestComp}" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=${out##*:}
    # Remove the directive
    out=${out%:*}
    if [ "${directive}" = "${out}" ]; then
        # There is not directive specified
        directive=0
    fi
    __prog_debug "${FUNCNAME[0]}: the completion directive is: ${directive}"
    __prog_debug "${FUNCNAME[0]}: the completions are: ${out[*]}"

    if [ $((directive & 1)) -ne 0 ]; then
        # Error code.  No completion.
        __prog_debug "${FUNCNAME[0]}: received error from custom completion go code"
        return
    else
        if [ $((directive & 2)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __prog_debug "${FUNCNAME[0]}: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & 4)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __prog_debug "${FUNCNAME[0]}: activating no file completion"
                compopt +o default
            fi
        fi

        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${out[*]}" -- "$cur")
    fi
}

__prog_handle_reply()
{
    __prog_debug "${FUNCNAME[0]}"
    local comp
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            while IFS='' read -r comp; do
                COMPREPLY+=("$comp")
            done < <(compgen -W "${allflags[*]}" -- "$cur")
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __prog_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __prog_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions=("${must_have_one_noun[@]}")
    elif [[ -n "${has_completion_function}" ]]; then
        # if a go completion function is provided, defer to that function
        completions=()
        __prog_handle_go_custom_completion
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "${completions[*]}" -- "$cur")

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${noun_aliases[*]}" -- "$cur")
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
		if declare -F __prog_custom_func >/dev/null; then
			# try command name qualified custom func
			__prog_custom_func
		else
			# otherwise fall back to unqualified for compatibility
			declare -F __custom_func >/dev/null && __custom_func
		fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__prog_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__prog_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
}

__prog_handle_flag()
{
    __prog_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __prog_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __prog_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __prog_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ ${words[c]} != *"="* ]] && __prog_contains_word "${words[c]}" "${two_word_flags[@]}"; then
			  __prog_debug "${FUNCNAME[0]}: found a flag ${words[c]}, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__prog_handle_noun()
{
    __prog_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __prog_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __prog_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__prog_handle_command()
{
    __prog_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_prog_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __prog_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__prog_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __prog_handle_reply
        return
    fi
    __prog_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __prog_handle_flag
    elif __prog_contains_word "${words[c]}" "${commands[@]}"; then
        __prog_handle_command
    elif [[ $c -eq 0 ]]; then
        __prog_handle_command
    elif __prog_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __prog_handle_command
        else
            __prog_handle_noun
        fi
    else
        __prog_handle_noun
    fi
    __prog_handle_word
}

_prog_completion_bash()
{
    last_command="prog_completion_bash"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--abool")
    flags+=("--abytes=")
    two_word_flags+=("--abytes")
    flags+=("--aduration=")
    two_word_flags+=("--aduration")
    flags+=("--afloat=")
    two_word_flags+=("--afloat")
    flags+=("--anint=")
    two_word_flags+=("--anint")
    flags+=("--anip=")
    two_word_flags+=("--anip")
    flags+=("--anipnet=")
    two_word_flags+=("--anipnet")
    flags+=("--astringslice=")
    two_word_flags+=("--astringslice")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--debug-address=")
    two_word_flags+=("--debug-address")
    flags+=("--debug-autostart")
    flags+=("--log-file-path=")
    two_word_flags+=("--log-file-path")
    flags+=("--log-format=")
    two_word_flags+=("--log-format")
    flags+=("--log-level=")
    two_word_flags+=("--log-level")
    flags+=("--log-target=")
    two_word_flags+=("--log-target")
    flags+=("--remote-first=")
    two_word_flags+=("--remote-first")
    flags+=("--remote-second=")
    two_word_flags+=("--remote-second")
    flags+=("--short=")
    two_word_flags+=("--short")
    two_word_flags+=("-s")
    flags+=("--store-another-sub-level=")
    two_word_flags+=("--store-another-sub-level")
    flags+=("--store-dir=")
    two_word_flags+=("--store-dir")
    two_word_flags+=("-S")
    flags+=("--store-log-level=")
    two_word_flags+=("--store-log-level")
    two_word_flags+=("-L")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_prog_completion_powershell()
{
    last_command="prog_completion_powershell"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--abool")
    flags+=("--abytes=")
    two_word_flags+=("--abytes")
    flags+=("--aduration=")
    two_word_flags+=("--aduration")
    flags+=("--afloat=")
    two_word_flags+=("--afloat")
    flags+=("--anint=")
    two_word_flags+=("--anint")
    flags+=("--anip=")
    two_word_flags+=("--anip")
    flags+=("--anipnet=")
    two_word_flags+=("--anipnet")
    flags+=("--astringslice=")
    two_word_flags+=("--astringslice")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--debug-address=")
    two_word_flags+=("--debug-address")
    flags+=("--debug-autostart")
    flags+=("--log-file-path=")
    two_word_flags+=("--log-file-path")
    flags+=("--log-format=")
    two_word_flags+=("--log-format")
    flags+=("--log-level=")
    two_word_flags+=("--log-level")
    flags+=("--log-target=")
    two_word_flags+=("--log-target")
    flags+=("--remote-first=")
    two_word_flags+=("--remote-first")
    flags+=("--remote-second=")
    two_word_flags+=("--remote-second")
    flags+=("--short=")
    two_word_flags+=("--short")
    two_word_flags+=("-s")
    flags+=("--store-another-sub-level=")
    two_word_flags+=("--store-another-sub-level")
    flags+=("--store-dir=")
    two_word_flags+=("--store-dir")
    two_word_flags+=("-S")
    flags+=("--store-log-level=")
    two_word_flags+=("--store-log-level")
    two_word_flags+=("-L")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_prog_completion_zsh()
{
    last_command="prog_completion_zsh"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--abool")
    flags+=("--abytes=")
    two_word_flags+=("--abytes")
    flags+=("--aduration=")
    two_word_flags+=("--aduration")
    flags+=("--afloat=")
    two_word_flags+=("--afloat")
    flags+=("--anint=")
    two_word_flags+=("--anint")
    flags+=("--anip=")
    two_word_flags+=("--anip")
    flags+=("--anipnet=")
    two_word_flags+=("--anipnet")
    flags+=("--astringslice=")
    two_word_flags+=("--astringslice")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--debug-address=")
    two_word_flags+=("--debug-address")
    flags+=("--debug-autostart")
    flags+=("--log-file-path=")
    two_word_flags+=("--log-file-path")
    flags+=("--log-format=")
    two_word_flags+=("--log-format")
    flags+=("--log-level=")
    two_word_flags+=("--log-level")
    flags+=("--log-target=")
    two_word_flags+=("--log-target")
    flags+=("--remote-first=")
    two_word_flags+=("--remote-first")
    flags+=("--remote-second=")
    two_word_flags+=("--remote-second")
    flags+=("--short=")
    two_word_flags+=("--short")
    two_word_flags+=("-s")
    flags+=("--store-another-sub-level=")
    two_word_flags+=("--store-another-sub-level")
    flags+=("--store-dir=")
    two_word_flags+=("--store-dir")
    two_word_flags+=("-S")
    flags+=("--store-log-level=")
    two_word_flags+=("--store-log-level")
    two_word_flags+=("-L")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_prog_completion()
{
    last_command="prog_completion"

    command_aliases=()

    commands=()
    commands+=("bash")
    commands+=("powershell")
    commands+=("zsh")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--abool")
    flags+=("--abytes=")
    two_word_flags+=("--abytes")
    flags+=("--aduration=")
    two_word_flags+=("--aduration")
    flags+=("--afloat=")
    two_word_flags+=("--afloat")
    flags+=("--anint=")
    two_word_flags+=("--anint")
    flags+=("--anip=")
    two_word_flags+=("--anip")
    flags+=("--anipnet=")
    two_word_flags+=("--anipnet")
    flags+=("--astringslice=")
    two_word_flags+=("--astringslice")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--debug-address=")
    two_word_flags+=("--debug-address")
    flags+=("--debug-autostart")
    flags+=("--log-file-path=")
    two_word_flags+=("--log-file-path")
    flags+=("--log-format=")
    two_word_flags+=("--log-format")
    flags+=("--log-level=")
    two_word_flags+=("--log-level")
    flags+=("--log-target=")
    two_word_flags+=("--log-target")
    flags+=("--remote-first=")
    two_word_flags+=("--remote-first")
    flags+=("--remote-second=")
    two_word_flags+=("--remote-second")
    flags+=("--short=")
    two_word_flags+=("--short")
    two_word_flags+=("-s")
    flags+=("--store-another-sub-level=")
    two_word_flags+=("--store-another-sub-level")
    flags+=("--store-dir=")
    two_word_flags+=("--store-dir")
    two_word_flags+=("-S")
    flags+=("--store-log-level=")
    two_word_flags+=("--store-log-level")
    two_word_flags+=("-L")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_prog_config_show()
{
    last_command="prog_config_show"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--abool")
    flags+=("--abytes=")
    two_word_flags+=("--abytes")
    flags+=("--aduration=")
    two_word_flags+=("--aduration")
    flags+=("--afloat=")
    two_word_flags+=("--afloat")
    flags+=("--anint=")
    two_word_flags+=("--anint")
    flags+=("--anip=")
    two_word_flags+=("--anip")
    flags+=("--anipnet=")
    two_word_flags+=("--anipnet")
    flags+=("--astringslice=")
    two_word_flags+=("--astringslice")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--debug-address=")
    two_word_flags+=("--debug-address")
    flags+=("--debug-autostart")
    flags+=("--log-file-path=")
    two_word_flags+=("--log-file-path")
    flags+=("--log-format=")
    two_word_flags+=("--log-format")
    flags+=("--log-level=")
    two_word_flags+=("--log-level")
    flags+=("--log-target=")
    two_word_flags+=("--log-target")
    flags+=("--remote-first=")
    two_word_flags+=("--remote-first")
    flags+=("--remote-second=")
    two_word_flags+=("--remote-second")
    flags+=("--short=")
    two_word_flags+=("--short")
    two_word_flags+=("-s")
    flags+=("--store-another-sub-level=")
    two_word_flags+=("--store-another-sub-level")
    flags+=("--store-dir=")
    two_word_flags+=("--store-dir")
    two_word_flags+=("-S")
    flags+=("--store-log-level=")
    two_word_flags+=("--store-log-level")
    two_word_flags+=("-L")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("dotenv")
    must_have_one_noun+=("env")
    must_have_one_noun+=("hcl")
    must_have_one_noun+=("ini")
    must_have_one_noun+=("json")
    must_have_one_noun+=("prop")
    must_have_one_noun+=("properties")
    must_have_one_noun+=("props")
    must_have_one_noun+=("toml")
    must_have_one_noun+=("yaml")
    must_have_one_noun+=("yml")
    noun_aliases=()
}

_prog_config()
{
    last_command="prog_config"

    command_aliases=()

    commands=()
    commands+=("show")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--abool")
    flags+=("--abytes=")
    two_word_flags+=("--abytes")
    flags+=("--aduration=")
    two_word_flags+=("--aduration")
    flags+=("--afloat=")
    two_word_flags+=("--afloat")
    flags+=("--anint=")
    two_word_flags+=("--anint")
    flags+=("--anip=")
    two_word_flags+=("--anip")
    flags+=("--anipnet=")
    two_word_flags+=("--anipnet")
    flags+=("--astringslice=")
    two_word_flags+=("--astringslice")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--debug-address=")
    two_word_flags+=("--debug-address")
    flags+=("--debug-autostart")
    flags+=("--log-file-path=")
    two_word_flags+=("--log-file-path")
    flags+=("--log-format=")
    two_word_flags+=("--log-format")
    flags+=("--log-level=")
    two_word_flags+=("--log-level")
    flags+=("--log-target=")
    two_word_flags+=("--log-target")
    flags+=("--remote-first=")
    two_word_flags+=("--remote-first")
    flags+=("--remote-second=")
    two_word_flags+=("--remote-second")
    flags+=("--short=")
    two_word_flags+=("--short")
    two_word_flags+=("-s")
    flags+=("--store-another-sub-level=")
    two_word_flags+=("--store-another-sub-level")
    flags+=("--store-dir=")
    two_word_flags+=("--store-dir")
    two_word_flags+=("-S")
    flags+=("--store-log-level=")
    two_word_flags+=("--store-log-level")
    two_word_flags+=("-L")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_prog_version()
{
    last_command="prog_version"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--abool")
    flags+=("--abytes=")
    two_word_flags+=("--abytes")
    flags+=("--aduration=")
    two_word_flags+=("--aduration")
    flags+=("--afloat=")
    two_word_flags+=("--afloat")
    flags+=("--anint=")
    two_word_flags+=("--anint")
    flags+=("--anip=")
    two_word_flags+=("--anip")
    flags+=("--anipnet=")
    two_word_flags+=("--anipnet")
    flags+=("--astringslice=")
    two_word_flags+=("--astringslice")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--debug-address=")
    two_word_flags+=("--debug-address")
    flags+=("--debug-autostart")
    flags+=("--log-file-path=")
    two_word_flags+=("--log-file-path")
    flags+=("--log-format=")
    two_word_flags+=("--log-format")
    flags+=("--log-level=")
    two_word_flags+=("--log-level")
    flags+=("--log-target=")
    two_word_flags+=("--log-target")
    flags+=("--remote-first=")
    two_word_flags+=("--remote-first")
    flags+=("--remote-second=")
    two_word_flags+=("--remote-second")
    flags+=("--short=")
    two_word_flags+=("--short")
    two_word_flags+=("-s")
    flags+=("--store-another-sub-level=")
    two_word_flags+=("--store-another-sub-level")
    flags+=("--store-dir=")
    two_word_flags+=("--store-dir")
    two_word_flags+=("-S")
    flags+=("--store-log-level=")
    two_word_flags+=("--store-log-level")
    two_word_flags+=("-L")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_prog_root_command()
{
    last_command="prog"

    command_aliases=()

    commands=()
    commands+=("completion")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("compl")
        aliashash["compl"]="completion"
    fi
    commands+=("config")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("conf")
        aliashash["conf"]="config"
    fi
    commands+=("version")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--abool")
    flags+=("--abytes=")
    two_word_flags+=("--abytes")
    flags+=("--aduration=")
    two_word_flags+=("--aduration")
    flags+=("--afloat=")
    two_word_flags+=("--afloat")
    flags+=("--anint=")
    two_word_flags+=("--anint")
    flags+=("--anip=")
    two_word_flags+=("--anip")
    flags+=("--anipnet=")
    two_word_flags+=("--anipnet")
    flags+=("--astringslice=")
    two_word_flags+=("--astringslice")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--debug-address=")
    two_word_flags+=("--debug-address")
    flags+=("--debug-autostart")
    flags+=("--log-file-path=")
    two_word_flags+=("--log-file-path")
    flags+=("--log-format=")
    two_word_flags+=("--log-format")
    flags+=("--log-level=")
    two_word_flags+=("--log-level")
    flags+=("--log-target=")
    two_word_flags+=("--log-target")
    flags+=("--remote-first=")
    two_word_flags+=("--remote-first")
    flags+=("--remote-second=")
    two_word_flags+=("--remote-second")
    flags+=("--short=")
    two_word_flags+=("--short")
    two_word_flags+=("-s")
    flags+=("--store-another-sub-level=")
    two_word_flags+=("--store-another-sub-level")
    flags+=("--store-dir=")
    two_word_flags+=("--store-dir")
    two_word_flags+=("-S")
    flags+=("--store-log-level=")
    two_word_flags+=("--store-log-level")
    two_word_flags+=("-L")
    flags+=("--version")
    local_nonpersistent_flags+=("--version")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_prog()
{
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __prog_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("prog")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function
    local last_command
    local nouns=()

    __prog_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_prog prog
else
    complete -o default -o nospace -F __start_prog prog
fi

# ex: ts=4 sw=4 et filetype=sh
