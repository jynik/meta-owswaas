DESCRIPTION = "Barebones image for Owen Wilson Saying 'Wow' as a Service"
LICENSE = "MIT"

inherit core-image

OWSWAAS_DIR ??= '/opt/owswaas'

IMAGE_INSTALL += "\
    kernel-modules \
    init-owswaas \ 
"

inherit extrausers
EXTRA_USERS_PARAMS = "\
    useradd -r -d ${OWSWAAS_DIR} -s /bin/false owswaas; \
"

