SUMMARY = "Default audio samples for Owen Wilson Saying 'Wow' as a Service"
LICENSE = "CLOSED"

SRC_URI = "https://www.dropbox.com/s/nlmsbsrqmrm3xxz/audio-samples.tar.xz?dl=1"
SRC_URI[md5sum] = "c7a7062e7874ad950738aa0b7d641bc4"
SRC_URI[sha256sum] = "31c568359696eb77be90d5da376e878287a0403fe659e24850cad6135d0a9f43" 

# Hack around the annoying suffix. Thanks dr0pb0x.
do_configure() {
    cd ${WORKDIR}
    if [ -f ${WORKDIR}/audio-samples.tar.xz?dl=1 ]; then
        mv ${WORKDIR}/audio-samples.tar.xz?dl=1 ${WORKDIR}/audio-samples.tar.xz
    fi
    tar xf ${WORKDIR}/audio-samples.tar.xz
}

do_install() {
    mkdir -p ${D}/audio/

    cp -r ${WORKDIR}/audio-samples/wowz ${D}/audio/wowz
    cp -r ${WORKDIR}/audio-samples/dlr ${D}/audio/dlr

    chmod -R -w ${D}/audio
    chmod -R +r ${D}/audio

    chmod 555 ${D}/audio
    chmod 555 ${D}/audio/wowz
    chmod 555 ${D}/audio/dlr
}

PACKAGES = "${PN}"
FILES_${PN} = "/audio"
