SUMMARY = "Override default (empty) asound.conf with one defaulting to card ${DEFAULT_ALSA_CARD}"

# Default is hifiberry-dac on Raspberry Pi
DEFAULT_ALSA_CARD ??= "0"

do_install_prepend() {
    echo "pcm.!default {"               >  ${WORKDIR}/asound.conf
    echo " type hw"                     >> ${WORKDIR}/asound.conf
    echo " card ${DEFAULT_ALSA_CARD}"   >> ${WORKDIR}/asound.conf
    echo "}"                            >> ${WORKDIR}/asound.conf
    echo ""                             >> ${WORKDIR}/asound.conf
    echo "ctl.!default {"               >> ${WORKDIR}/asound.conf
    echo " card ${DEFAULT_ALSA_CARD}"   >> ${WORKDIR}/asound.conf
    echo "}"                            >> ${WORKDIR}/asound.conf
}
