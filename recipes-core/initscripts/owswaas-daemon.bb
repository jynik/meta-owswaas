DESCRIPTION = "Startup script for owswaas"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

inherit useradd

RDEPENDS_${PN} = "owswaas daemonize audio-samples"

SRC_URI = "file://owswaas-daemon"

S = "${WORKDIR}"

do_configure() {
    if [ ! -z "${GPIO_BTN}" ]; then
        sed -i -e "s/GPIO_BTN=[0-9]+/GPIO_BTN=${GPIO_BTN}/" ${S}/owswaas-daemon
    fi

    if [ ! -z "${GPIO_SW}" ]; then
        sed -i -e "s/GPIO_SW=[0-9]+/GPIO_SW=${GPIO_SW}/" ${S}/owswaas-daemon

    fi

    if [ ! -z "${GPIO_LED}" ]; then
        sed -i -e "s/GPIO_LED=[0-9]+/GPIO_LED=${GPIO_LED}/" ${S}/owswaas-daemon

    fi
}

do_install() {
    install -d ${D}/${sysconfdir}/init.d
    install -d ${D}/${sysconfdir}/rcS.d
    install -d ${D}/${sysconfdir}/rc1.d
    install -d ${D}/${sysconfdir}/rc2.d
    install -d ${D}/${sysconfdir}/rc3.d
    install -d ${D}/${sysconfdir}/rc4.d
    install -d ${D}/${sysconfdir}/rc5.d

    install -m 0755 ${S}/owswaas-daemon ${D}/${sysconfdir}/init.d/owswaas-daemon

    ln -sf ../init.d/owswaas-daemon ${D}${sysconfdir}/rc5.d/S90-owswaas
    ln -sf ../init.d/owswaas-daemon ${D}${sysconfdir}/rc1.d/K90-owswaas
    ln -sf ../init.d/owswaas-daemon ${D}${sysconfdir}/rc2.d/K90-owswaas
    ln -sf ../init.d/owswaas-daemon ${D}${sysconfdir}/rc3.d/K90-owswaas
    ln -sf ../init.d/owswaas-daemon ${D}${sysconfdir}/rc4.d/K90-owswaas
    ln -sf ../init.d/owswaas-daemon ${D}${sysconfdir}/rc5.d/K90-owswaas
}

USERADD_PACKAGES = "${PN}"
USERADD_PARAM_${PN} = "-r -s /bin/false owswaas"
GROUPMEMS_PARAM_${PN} = "-a owswaas -g audio"
