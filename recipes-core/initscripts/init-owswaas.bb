DESCRIPTION = "Startup script for owswaas"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

RDEPENDS_${PN} = "owswaas daemonize"

SRC_URI = "file://init-owswaas"

S = "${WORKDIR}"

BTN_GPIO_raspberrypi0 ??= "24"

do_configure() {
    sed -e "s/#BTN_GPIO#/${BTN_GPIO}/" ${S}/init-owswaas > ${S}/owswaas
}

do_install() {
    install -d ${D}/${sysconfdir}/init.d
    install -d ${D}/${sysconfdir}/rcS.d
    install -d ${D}/${sysconfdir}/rc1.d
    install -d ${D}/${sysconfdir}/rc2.d
    install -d ${D}/${sysconfdir}/rc3.d
    install -d ${D}/${sysconfdir}/rc4.d
    install -d ${D}/${sysconfdir}/rc5.d

    install -m 0755 ${S}/owswaas ${D}/${sysconfdir}/init.d/owswaas

    ln -sf ../init.d/owswaas ${D}${sysconfdir}/rc5.d/S90-owswaas
    ln -sf ../init.d/owswaas ${D}${sysconfdir}/rc1.d/K90-owswaas
    ln -sf ../init.d/owswaas ${D}${sysconfdir}/rc2.d/K90-owswaas
    ln -sf ../init.d/owswaas ${D}${sysconfdir}/rc3.d/K90-owswaas
    ln -sf ../init.d/owswaas ${D}${sysconfdir}/rc4.d/K90-owswaas
    ln -sf ../init.d/owswaas ${D}${sysconfdir}/rc5.d/K90-owswaas
}
