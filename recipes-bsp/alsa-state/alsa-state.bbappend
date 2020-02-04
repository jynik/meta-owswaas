SUMMARY = "Override default (empty) asound.conf with one defaulting to card ${DEFAULT_ALSA_CARD}"

# Default is USB Speaker @ 0
DEFAULT_ALSA_CARD ??= "0"

do_install_prepend() {
    echo "pcm.usb_speaker {"            >  ${WORKDIR}/asound.conf
    echo " type hw"                     >> ${WORKDIR}/asound.conf
    echo " card ${DEFAULT_ALSA_CARD}"   >> ${WORKDIR}/asound.conf
    echo "}"                            >> ${WORKDIR}/asound.conf
    echo ""                             >> ${WORKDIR}/asound.conf
    echo "ctl.usb_speaker {"            >> ${WORKDIR}/asound.conf
    echo " card ${DEFAULT_ALSA_CARD}"   >> ${WORKDIR}/asound.conf
    echo "}"                            >> ${WORKDIR}/asound.conf
}
