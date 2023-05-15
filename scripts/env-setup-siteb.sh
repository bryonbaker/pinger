#! /bin/bash

export KUBECONFIG=$HOME/.kube/config-siteb
export PS1="SITE B: \[$(tput setaf 2)\]\u@\h\[$(tput setaf 7)\]:\$(parse_git_branch)\[$(tput setaf 6)\]\W\\$ \[$(tput setaf 7)\]\[$(tput sgr0)\]"
