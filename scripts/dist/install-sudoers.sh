#!/usr/bin/env bash

# Installs the sudoers drop-in that lets the entrypoint run the init script as root.
# Pass --develop in development images, where all targets and scripts may be run with sudo.
# bash <(curl -s https://raw.githubusercontent.com/photoprism/photoprism/develop/scripts/dist/install-sudoers.sh)

PATH="/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin:/scripts:$PATH"

# Abort if not executed as root.
if [[ $(id -u) != "0" ]]; then
  echo "Usage: run ${0##*/} as root" 1>&2
  exit 1
fi

set -e

INIT_SCRIPT="/scripts/entrypoint-init.sh"
SUDOERS_FILE="/etc/sudoers.d/init"

# Variables the init script and the targets it runs read from the environment. Only these are
# passed on, so that the caller cannot supply the ones that change how a command interprets
# its input, such as the variables GNU make accepts options and additional makefiles through.
# The proxy variables are included because the init targets download packages and models, and
# nothing else provides them; every download uses TLS, and apt verifies its own signatures.
INIT_ENV="DOCKER_ENV DOCKER_TAG BUILD_ARCH TF_DRIVER TF_VERSION ONNX_GPU ONNX_VERSION DEBIAN_FRONTEND \
http_proxy https_proxy ftp_proxy all_proxy no_proxy HTTP_PROXY HTTPS_PROXY FTP_PROXY ALL_PROXY NO_PROXY \
PHOTOPRISM_*"

if [[ $1 == "--develop" ]]; then
  # Development images allow every target and script to be run with sudo.
  SUDOERS_FILE="/etc/sudoers.d/all"
  printf '%s\n' \
    "Defaults env_keep += \"${INIT_ENV}\"" \
    'ALL ALL=(ALL) NOPASSWD:SETENV: ALL' \
    > "${SUDOERS_FILE}"
else
  printf '%s\n' \
    "Cmnd_Alias PHOTOPRISM_INIT_CMND = ${INIT_SCRIPT} \"\"" \
    "Defaults!PHOTOPRISM_INIT_CMND env_keep += \"${INIT_ENV}\"" \
    'ALL ALL=(root) NOPASSWD: PHOTOPRISM_INIT_CMND' \
    > "${SUDOERS_FILE}"
fi

chmod 0440 "${SUDOERS_FILE}"

# Reject a malformed file at build time, as sudo would otherwise ignore it at runtime.
visudo -c -f "${SUDOERS_FILE}"

echo "✅ Installed ${SUDOERS_FILE}."
