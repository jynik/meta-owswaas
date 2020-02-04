DESCRIPTION = "Owen Wilson Saying 'Wow' as a Service"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

RDEPENDS_${PN} = "audio-samples alsa-utils-aplay alsa-utils-amixer"

inherit go
inherit godep

GO_IMPORT = "owswaas"

SRC_URI = "\
    file://owswaas.go \
"

S = "${WORKDIR}"

# I really don't want to host an entire package repo for this ridiculous thing.
# Just fake it to make it.
do_configure_prepend() {
    mkdir -p ${S}/src/owswaas
    cp ${S}/owswaas.go ${S}/src/owswaas
}

# Gross workaround for...
#
# ERROR: owswaas-1.0-r0 do_package_qa: QA Issue: 
#   /usr/lib/go/pkg/dep/sources/https---go.googlesource.com-sys/windows/mkknownfolderids.bash
#   contained in package owswaas-staticdev requires /bin/bash, but no providers found in 
#   RDEPENDS_owswaas-staticdev? [file-rdeps]
#
# It sure does... but I'm also not building anything for Windows so let's just
# burn this with fire before the package QA finds it.
#
do_compile_append() {
    rm -r ${B}/pkg/dep/sources/https---go.googlesource.com-sys/windows
}

FILES_${PN} += "${bindir}/owswaas"
