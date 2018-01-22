DESCRIPTION = "Owen Wilson Saying 'Wow' as a Service"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

SRC_URI = "file://owswaas"

RDEPENDS_${PN} = "wowz alsa-utils-aplay"

S = "${WORKDIR}"

do_install() {
    install -d ${D}/${sbindir}
    install -m 0744 ${S}/owswaas ${D}${sbindir}/owswaas
}
