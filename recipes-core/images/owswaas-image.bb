DESCRIPTION = "Barebones image for Owen Wilson Saying 'Wow' as a Service"
LICENSE = "MIT"

inherit core-image

IMAGE_INSTALL +=  "\
    kernel-modules \
    owswaas-daemon \
"
