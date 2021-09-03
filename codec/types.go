package codec

//#include <libavformat/avformat.h>
import "C"

const (
	EAGAIN = C.EAGAIN
)

// CODEC ID
const (
	AV_CODEC_ID_NONE               = C.AV_CODEC_ID_NONE
	AV_CODEC_ID_MPEG1VIDEO         = C.AV_CODEC_ID_MPEG1VIDEO
	AV_CODEC_ID_H261               = C.AV_CODEC_ID_H261
	AV_CODEC_ID_H263               = C.AV_CODEC_ID_H263
	AV_CODEC_ID_RV10               = C.AV_CODEC_ID_RV10
	AV_CODEC_ID_RV20               = C.AV_CODEC_ID_RV20
	AV_CODEC_ID_MJPEG              = C.AV_CODEC_ID_MJPEG
	AV_CODEC_ID_MJPEGB             = C.AV_CODEC_ID_MJPEGB
	AV_CODEC_ID_LJPEG              = C.AV_CODEC_ID_LJPEG
	AV_CODEC_ID_SP5X               = C.AV_CODEC_ID_SP5X
	AV_CODEC_ID_JPEGLS             = C.AV_CODEC_ID_JPEGLS
	AV_CODEC_ID_MPEG4              = C.AV_CODEC_ID_MPEG4
	AV_CODEC_ID_RAWVIDEO           = C.AV_CODEC_ID_RAWVIDEO
	AV_CODEC_ID_MSMPEG4V1          = C.AV_CODEC_ID_MSMPEG4V1
	AV_CODEC_ID_MSMPEG4V2          = C.AV_CODEC_ID_MSMPEG4V2
	AV_CODEC_ID_MSMPEG4V3          = C.AV_CODEC_ID_MSMPEG4V3
	AV_CODEC_ID_WMV1               = C.AV_CODEC_ID_WMV1
	AV_CODEC_ID_WMV2               = C.AV_CODEC_ID_WMV2
	AV_CODEC_ID_H263P              = C.AV_CODEC_ID_H263P
	AV_CODEC_ID_H263I              = C.AV_CODEC_ID_H263I
	AV_CODEC_ID_FLV1               = C.AV_CODEC_ID_FLV1
	AV_CODEC_ID_SVQ1               = C.AV_CODEC_ID_SVQ1
	AV_CODEC_ID_SVQ3               = C.AV_CODEC_ID_SVQ3
	AV_CODEC_ID_DVVIDEO            = C.AV_CODEC_ID_DVVIDEO
	AV_CODEC_ID_HUFFYUV            = C.AV_CODEC_ID_HUFFYUV
	AV_CODEC_ID_CYUV               = C.AV_CODEC_ID_CYUV
	AV_CODEC_ID_H264               = C.AV_CODEC_ID_H264
	AV_CODEC_ID_INDEO3             = C.AV_CODEC_ID_INDEO3
	AV_CODEC_ID_VP3                = C.AV_CODEC_ID_VP3
	AV_CODEC_ID_THEORA             = C.AV_CODEC_ID_THEORA
	AV_CODEC_ID_ASV1               = C.AV_CODEC_ID_ASV1
	AV_CODEC_ID_ASV2               = C.AV_CODEC_ID_ASV2
	AV_CODEC_ID_FFV1               = C.AV_CODEC_ID_FFV1
	AV_CODEC_ID_4XM                = C.AV_CODEC_ID_4XM
	AV_CODEC_ID_VCR1               = C.AV_CODEC_ID_VCR1
	AV_CODEC_ID_CLJR               = C.AV_CODEC_ID_CLJR
	AV_CODEC_ID_MDEC               = C.AV_CODEC_ID_MDEC
	AV_CODEC_ID_ROQ                = C.AV_CODEC_ID_ROQ
	AV_CODEC_ID_INTERPLAY_VIDEO    = C.AV_CODEC_ID_INTERPLAY_VIDEO
	AV_CODEC_ID_XAN_WC3            = C.AV_CODEC_ID_XAN_WC3
	AV_CODEC_ID_XAN_WC4            = C.AV_CODEC_ID_XAN_WC4
	AV_CODEC_ID_RPZA               = C.AV_CODEC_ID_RPZA
	AV_CODEC_ID_CINEPAK            = C.AV_CODEC_ID_CINEPAK
	AV_CODEC_ID_WS_VQA             = C.AV_CODEC_ID_WS_VQA
	AV_CODEC_ID_MSRLE              = C.AV_CODEC_ID_MSRLE
	AV_CODEC_ID_MSVIDEO1           = C.AV_CODEC_ID_MSVIDEO1
	AV_CODEC_ID_IDCIN              = C.AV_CODEC_ID_IDCIN
	AV_CODEC_ID_8BPS               = C.AV_CODEC_ID_8BPS
	AV_CODEC_ID_SMC                = C.AV_CODEC_ID_SMC
	AV_CODEC_ID_FLIC               = C.AV_CODEC_ID_FLIC
	AV_CODEC_ID_TRUEMOTION1        = C.AV_CODEC_ID_TRUEMOTION1
	AV_CODEC_ID_VMDVIDEO           = C.AV_CODEC_ID_VMDVIDEO
	AV_CODEC_ID_MSZH               = C.AV_CODEC_ID_MSZH
	AV_CODEC_ID_ZLIB               = C.AV_CODEC_ID_ZLIB
	AV_CODEC_ID_QTRLE              = C.AV_CODEC_ID_QTRLE
	AV_CODEC_ID_TSCC               = C.AV_CODEC_ID_TSCC
	AV_CODEC_ID_ULTI               = C.AV_CODEC_ID_ULTI
	AV_CODEC_ID_QDRAW              = C.AV_CODEC_ID_QDRAW
	AV_CODEC_ID_VIXL               = C.AV_CODEC_ID_VIXL
	AV_CODEC_ID_QPEG               = C.AV_CODEC_ID_QPEG
	AV_CODEC_ID_PNG                = C.AV_CODEC_ID_PNG
	AV_CODEC_ID_PPM                = C.AV_CODEC_ID_PPM
	AV_CODEC_ID_PBM                = C.AV_CODEC_ID_PBM
	AV_CODEC_ID_PGM                = C.AV_CODEC_ID_PGM
	AV_CODEC_ID_PGMYUV             = C.AV_CODEC_ID_PGMYUV
	AV_CODEC_ID_PAM                = C.AV_CODEC_ID_PAM
	AV_CODEC_ID_FFVHUFF            = C.AV_CODEC_ID_FFVHUFF
	AV_CODEC_ID_RV30               = C.AV_CODEC_ID_RV30
	AV_CODEC_ID_RV40               = C.AV_CODEC_ID_RV40
	AV_CODEC_ID_VC1                = C.AV_CODEC_ID_VC1
	AV_CODEC_ID_WMV3               = C.AV_CODEC_ID_WMV3
	AV_CODEC_ID_LOCO               = C.AV_CODEC_ID_LOCO
	AV_CODEC_ID_WNV1               = C.AV_CODEC_ID_WNV1
	AV_CODEC_ID_AASC               = C.AV_CODEC_ID_AASC
	AV_CODEC_ID_INDEO2             = C.AV_CODEC_ID_INDEO2
	AV_CODEC_ID_FRAPS              = C.AV_CODEC_ID_FRAPS
	AV_CODEC_ID_TRUEMOTION2        = C.AV_CODEC_ID_TRUEMOTION2
	AV_CODEC_ID_BMP                = C.AV_CODEC_ID_BMP
	AV_CODEC_ID_CSCD               = C.AV_CODEC_ID_CSCD
	AV_CODEC_ID_MMVIDEO            = C.AV_CODEC_ID_MMVIDEO
	AV_CODEC_ID_ZMBV               = C.AV_CODEC_ID_ZMBV
	AV_CODEC_ID_AVS                = C.AV_CODEC_ID_AVS
	AV_CODEC_ID_SMACKVIDEO         = C.AV_CODEC_ID_SMACKVIDEO
	AV_CODEC_ID_NUV                = C.AV_CODEC_ID_NUV
	AV_CODEC_ID_KMVC               = C.AV_CODEC_ID_KMVC
	AV_CODEC_ID_FLASHSV            = C.AV_CODEC_ID_FLASHSV
	AV_CODEC_ID_CAVS               = C.AV_CODEC_ID_CAVS
	AV_CODEC_ID_JPEG2000           = C.AV_CODEC_ID_JPEG2000
	AV_CODEC_ID_VMNC               = C.AV_CODEC_ID_VMNC
	AV_CODEC_ID_VP5                = C.AV_CODEC_ID_VP5
	AV_CODEC_ID_VP6                = C.AV_CODEC_ID_VP6
	AV_CODEC_ID_VP6F               = C.AV_CODEC_ID_VP6F
	AV_CODEC_ID_TARGA              = C.AV_CODEC_ID_TARGA
	AV_CODEC_ID_DSICINVIDEO        = C.AV_CODEC_ID_DSICINVIDEO
	AV_CODEC_ID_TIERTEXSEQVIDEO    = C.AV_CODEC_ID_TIERTEXSEQVIDEO
	AV_CODEC_ID_TIFF               = C.AV_CODEC_ID_TIFF
	AV_CODEC_ID_GIF                = C.AV_CODEC_ID_GIF
	AV_CODEC_ID_DXA                = C.AV_CODEC_ID_DXA
	AV_CODEC_ID_DNXHD              = C.AV_CODEC_ID_DNXHD
	AV_CODEC_ID_THP                = C.AV_CODEC_ID_THP
	AV_CODEC_ID_SGI                = C.AV_CODEC_ID_SGI
	AV_CODEC_ID_C93                = C.AV_CODEC_ID_C93
	AV_CODEC_ID_BETHSOFTVID        = C.AV_CODEC_ID_BETHSOFTVID
	AV_CODEC_ID_PTX                = C.AV_CODEC_ID_PTX
	AV_CODEC_ID_TXD                = C.AV_CODEC_ID_TXD
	AV_CODEC_ID_VP6A               = C.AV_CODEC_ID_VP6A
	AV_CODEC_ID_AMV                = C.AV_CODEC_ID_AMV
	AV_CODEC_ID_VB                 = C.AV_CODEC_ID_VB
	AV_CODEC_ID_PCX                = C.AV_CODEC_ID_PCX
	AV_CODEC_ID_SUNRAST            = C.AV_CODEC_ID_SUNRAST
	AV_CODEC_ID_INDEO4             = C.AV_CODEC_ID_INDEO4
	AV_CODEC_ID_INDEO5             = C.AV_CODEC_ID_INDEO5
	AV_CODEC_ID_MIMIC              = C.AV_CODEC_ID_MIMIC
	AV_CODEC_ID_RL2                = C.AV_CODEC_ID_RL2
	AV_CODEC_ID_ESCAPE124          = C.AV_CODEC_ID_ESCAPE124
	AV_CODEC_ID_DIRAC              = C.AV_CODEC_ID_DIRAC
	AV_CODEC_ID_BFI                = C.AV_CODEC_ID_BFI
	AV_CODEC_ID_CMV                = C.AV_CODEC_ID_CMV
	AV_CODEC_ID_MOTIONPIXELS       = C.AV_CODEC_ID_MOTIONPIXELS
	AV_CODEC_ID_TGV                = C.AV_CODEC_ID_TGV
	AV_CODEC_ID_TGQ                = C.AV_CODEC_ID_TGQ
	AV_CODEC_ID_TQI                = C.AV_CODEC_ID_TQI
	AV_CODEC_ID_AURA               = C.AV_CODEC_ID_AURA
	AV_CODEC_ID_AURA2              = C.AV_CODEC_ID_AURA2
	AV_CODEC_ID_V210X              = C.AV_CODEC_ID_V210X
	AV_CODEC_ID_TMV                = C.AV_CODEC_ID_TMV
	AV_CODEC_ID_V210               = C.AV_CODEC_ID_V210
	AV_CODEC_ID_DPX                = C.AV_CODEC_ID_DPX
	AV_CODEC_ID_MAD                = C.AV_CODEC_ID_MAD
	AV_CODEC_ID_FRWU               = C.AV_CODEC_ID_FRWU
	AV_CODEC_ID_FLASHSV2           = C.AV_CODEC_ID_FLASHSV2
	AV_CODEC_ID_CDGRAPHICS         = C.AV_CODEC_ID_CDGRAPHICS
	AV_CODEC_ID_R210               = C.AV_CODEC_ID_R210
	AV_CODEC_ID_ANM                = C.AV_CODEC_ID_ANM
	AV_CODEC_ID_BINKVIDEO          = C.AV_CODEC_ID_BINKVIDEO
	AV_CODEC_ID_IFF_ILBM           = C.AV_CODEC_ID_IFF_ILBM
	AV_CODEC_ID_KGV1               = C.AV_CODEC_ID_KGV1
	AV_CODEC_ID_YOP                = C.AV_CODEC_ID_YOP
	AV_CODEC_ID_VP8                = C.AV_CODEC_ID_VP8
	AV_CODEC_ID_PICTOR             = C.AV_CODEC_ID_PICTOR
	AV_CODEC_ID_ANSI               = C.AV_CODEC_ID_ANSI
	AV_CODEC_ID_A64_MULTI          = C.AV_CODEC_ID_A64_MULTI
	AV_CODEC_ID_A64_MULTI5         = C.AV_CODEC_ID_A64_MULTI5
	AV_CODEC_ID_R10K               = C.AV_CODEC_ID_R10K
	AV_CODEC_ID_MXPEG              = C.AV_CODEC_ID_MXPEG
	AV_CODEC_ID_LAGARITH           = C.AV_CODEC_ID_LAGARITH
	AV_CODEC_ID_PRORES             = C.AV_CODEC_ID_PRORES
	AV_CODEC_ID_JV                 = C.AV_CODEC_ID_JV
	AV_CODEC_ID_DFA                = C.AV_CODEC_ID_DFA
	AV_CODEC_ID_WMV3IMAGE          = C.AV_CODEC_ID_WMV3IMAGE
	AV_CODEC_ID_VC1IMAGE           = C.AV_CODEC_ID_VC1IMAGE
	AV_CODEC_ID_UTVIDEO            = C.AV_CODEC_ID_UTVIDEO
	AV_CODEC_ID_BMV_VIDEO          = C.AV_CODEC_ID_BMV_VIDEO
	AV_CODEC_ID_VBLE               = C.AV_CODEC_ID_VBLE
	AV_CODEC_ID_DXTORY             = C.AV_CODEC_ID_DXTORY
	AV_CODEC_ID_V410               = C.AV_CODEC_ID_V410
	AV_CODEC_ID_XWD                = C.AV_CODEC_ID_XWD
	AV_CODEC_ID_CDXL               = C.AV_CODEC_ID_CDXL
	AV_CODEC_ID_XBM                = C.AV_CODEC_ID_XBM
	AV_CODEC_ID_ZEROCODEC          = C.AV_CODEC_ID_ZEROCODEC
	AV_CODEC_ID_MSS1               = C.AV_CODEC_ID_MSS1
	AV_CODEC_ID_MSA1               = C.AV_CODEC_ID_MSA1
	AV_CODEC_ID_TSCC2              = C.AV_CODEC_ID_TSCC2
	AV_CODEC_ID_MTS2               = C.AV_CODEC_ID_MTS2
	AV_CODEC_ID_CLLC               = C.AV_CODEC_ID_CLLC
	AV_CODEC_ID_MSS2               = C.AV_CODEC_ID_MSS2
	AV_CODEC_ID_VP9                = C.AV_CODEC_ID_VP9
	AV_CODEC_ID_AIC                = C.AV_CODEC_ID_AIC
	AV_CODEC_ID_ESCAPE130          = C.AV_CODEC_ID_ESCAPE130
	AV_CODEC_ID_G2M                = C.AV_CODEC_ID_G2M
	AV_CODEC_ID_WEBP               = C.AV_CODEC_ID_WEBP
	AV_CODEC_ID_HNM4_VIDEO         = C.AV_CODEC_ID_HNM4_VIDEO
	AV_CODEC_ID_HEVC               = C.AV_CODEC_ID_HEVC
	AV_CODEC_ID_FIC                = C.AV_CODEC_ID_FIC
	AV_CODEC_ID_ALIAS_PIX          = C.AV_CODEC_ID_ALIAS_PIX
	AV_CODEC_ID_BRENDER_PIX        = C.AV_CODEC_ID_BRENDER_PIX
	AV_CODEC_ID_PAF_VIDEO          = C.AV_CODEC_ID_PAF_VIDEO
	AV_CODEC_ID_EXR                = C.AV_CODEC_ID_EXR
	AV_CODEC_ID_VP7                = C.AV_CODEC_ID_VP7
	AV_CODEC_ID_SANM               = C.AV_CODEC_ID_SANM
	AV_CODEC_ID_SGIRLE             = C.AV_CODEC_ID_SGIRLE
	AV_CODEC_ID_MVC1               = C.AV_CODEC_ID_MVC1
	AV_CODEC_ID_MVC2               = C.AV_CODEC_ID_MVC2
	AV_CODEC_ID_HQX                = C.AV_CODEC_ID_HQX
	AV_CODEC_ID_TDSC               = C.AV_CODEC_ID_TDSC
	AV_CODEC_ID_HQ_HQA             = C.AV_CODEC_ID_HQ_HQA
	AV_CODEC_ID_HAP                = C.AV_CODEC_ID_HAP
	AV_CODEC_ID_DDS                = C.AV_CODEC_ID_DDS
	AV_CODEC_ID_DXV                = C.AV_CODEC_ID_DXV
	AV_CODEC_ID_SCREENPRESSO       = C.AV_CODEC_ID_SCREENPRESSO
	AV_CODEC_ID_RSCC               = C.AV_CODEC_ID_RSCC
	AV_CODEC_ID_AVS2               = C.AV_CODEC_ID_AVS2
	AV_CODEC_ID_Y41P               = C.AV_CODEC_ID_Y41P
	AV_CODEC_ID_AVRP               = C.AV_CODEC_ID_AVRP
	AV_CODEC_ID_012V               = C.AV_CODEC_ID_012V
	AV_CODEC_ID_AVUI               = C.AV_CODEC_ID_AVUI
	AV_CODEC_ID_AYUV               = C.AV_CODEC_ID_AYUV
	AV_CODEC_ID_TARGA_Y216         = C.AV_CODEC_ID_TARGA_Y216
	AV_CODEC_ID_V308               = C.AV_CODEC_ID_V308
	AV_CODEC_ID_V408               = C.AV_CODEC_ID_V408
	AV_CODEC_ID_YUV4               = C.AV_CODEC_ID_YUV4
	AV_CODEC_ID_AVRN               = C.AV_CODEC_ID_AVRN
	AV_CODEC_ID_CPIA               = C.AV_CODEC_ID_CPIA
	AV_CODEC_ID_XFACE              = C.AV_CODEC_ID_XFACE
	AV_CODEC_ID_SNOW               = C.AV_CODEC_ID_SNOW
	AV_CODEC_ID_SMVJPEG            = C.AV_CODEC_ID_SMVJPEG
	AV_CODEC_ID_APNG               = C.AV_CODEC_ID_APNG
	AV_CODEC_ID_DAALA              = C.AV_CODEC_ID_DAALA
	AV_CODEC_ID_CFHD               = C.AV_CODEC_ID_CFHD
	AV_CODEC_ID_TRUEMOTION2RT      = C.AV_CODEC_ID_TRUEMOTION2RT
	AV_CODEC_ID_M101               = C.AV_CODEC_ID_M101
	AV_CODEC_ID_MAGICYUV           = C.AV_CODEC_ID_MAGICYUV
	AV_CODEC_ID_SHEERVIDEO         = C.AV_CODEC_ID_SHEERVIDEO
	AV_CODEC_ID_YLC                = C.AV_CODEC_ID_YLC
	AV_CODEC_ID_PSD                = C.AV_CODEC_ID_PSD
	AV_CODEC_ID_PIXLET             = C.AV_CODEC_ID_PIXLET
	AV_CODEC_ID_SPEEDHQ            = C.AV_CODEC_ID_SPEEDHQ
	AV_CODEC_ID_FMVC               = C.AV_CODEC_ID_FMVC
	AV_CODEC_ID_SCPR               = C.AV_CODEC_ID_SCPR
	AV_CODEC_ID_CLEARVIDEO         = C.AV_CODEC_ID_CLEARVIDEO
	AV_CODEC_ID_XPM                = C.AV_CODEC_ID_XPM
	AV_CODEC_ID_AV1                = C.AV_CODEC_ID_AV1
	AV_CODEC_ID_BITPACKED          = C.AV_CODEC_ID_BITPACKED
	AV_CODEC_ID_MSCC               = C.AV_CODEC_ID_MSCC
	AV_CODEC_ID_SRGC               = C.AV_CODEC_ID_SRGC
	AV_CODEC_ID_SVG                = C.AV_CODEC_ID_SVG
	AV_CODEC_ID_GDV                = C.AV_CODEC_ID_GDV
	AV_CODEC_ID_FITS               = C.AV_CODEC_ID_FITS
	AV_CODEC_ID_IMM4               = C.AV_CODEC_ID_IMM4
	AV_CODEC_ID_PROSUMER           = C.AV_CODEC_ID_PROSUMER
	AV_CODEC_ID_MWSC               = C.AV_CODEC_ID_MWSC
	AV_CODEC_ID_WCMV               = C.AV_CODEC_ID_WCMV
	AV_CODEC_ID_RASC               = C.AV_CODEC_ID_RASC
	AV_CODEC_ID_FIRST_AUDIO        = C.AV_CODEC_ID_FIRST_AUDIO
	AV_CODEC_ID_PCM_S16LE          = C.AV_CODEC_ID_PCM_S16LE
	AV_CODEC_ID_PCM_S16BE          = C.AV_CODEC_ID_PCM_S16BE
	AV_CODEC_ID_PCM_U16LE          = C.AV_CODEC_ID_PCM_U16LE
	AV_CODEC_ID_PCM_U16BE          = C.AV_CODEC_ID_PCM_U16BE
	AV_CODEC_ID_PCM_S8             = C.AV_CODEC_ID_PCM_S8
	AV_CODEC_ID_PCM_U8             = C.AV_CODEC_ID_PCM_U8
	AV_CODEC_ID_PCM_MULAW          = C.AV_CODEC_ID_PCM_MULAW
	AV_CODEC_ID_PCM_ALAW           = C.AV_CODEC_ID_PCM_ALAW
	AV_CODEC_ID_PCM_S32LE          = C.AV_CODEC_ID_PCM_S32LE
	AV_CODEC_ID_PCM_S32BE          = C.AV_CODEC_ID_PCM_S32BE
	AV_CODEC_ID_PCM_U32LE          = C.AV_CODEC_ID_PCM_U32LE
	AV_CODEC_ID_PCM_U32BE          = C.AV_CODEC_ID_PCM_U32BE
	AV_CODEC_ID_PCM_S24LE          = C.AV_CODEC_ID_PCM_S24LE
	AV_CODEC_ID_PCM_S24BE          = C.AV_CODEC_ID_PCM_S24BE
	AV_CODEC_ID_PCM_U24LE          = C.AV_CODEC_ID_PCM_U24LE
	AV_CODEC_ID_PCM_U24BE          = C.AV_CODEC_ID_PCM_U24BE
	AV_CODEC_ID_PCM_S24DAUD        = C.AV_CODEC_ID_PCM_S24DAUD
	AV_CODEC_ID_PCM_ZORK           = C.AV_CODEC_ID_PCM_ZORK
	AV_CODEC_ID_PCM_S16LE_PLANAR   = C.AV_CODEC_ID_PCM_S16LE_PLANAR
	AV_CODEC_ID_PCM_DVD            = C.AV_CODEC_ID_PCM_DVD
	AV_CODEC_ID_PCM_F32BE          = C.AV_CODEC_ID_PCM_F32BE
	AV_CODEC_ID_PCM_F32LE          = C.AV_CODEC_ID_PCM_F32LE
	AV_CODEC_ID_PCM_F64BE          = C.AV_CODEC_ID_PCM_F64BE
	AV_CODEC_ID_PCM_F64LE          = C.AV_CODEC_ID_PCM_F64LE
	AV_CODEC_ID_PCM_BLURAY         = C.AV_CODEC_ID_PCM_BLURAY
	AV_CODEC_ID_PCM_LXF            = C.AV_CODEC_ID_PCM_LXF
	AV_CODEC_ID_S302M              = C.AV_CODEC_ID_S302M
	AV_CODEC_ID_PCM_S8_PLANAR      = C.AV_CODEC_ID_PCM_S8_PLANAR
	AV_CODEC_ID_PCM_S24LE_PLANAR   = C.AV_CODEC_ID_PCM_S24LE_PLANAR
	AV_CODEC_ID_PCM_S32LE_PLANAR   = C.AV_CODEC_ID_PCM_S32LE_PLANAR
	AV_CODEC_ID_PCM_S16BE_PLANAR   = C.AV_CODEC_ID_PCM_S16BE_PLANAR
	AV_CODEC_ID_PCM_S64LE          = C.AV_CODEC_ID_PCM_S64LE
	AV_CODEC_ID_PCM_S64BE          = C.AV_CODEC_ID_PCM_S64BE
	AV_CODEC_ID_PCM_F16LE          = C.AV_CODEC_ID_PCM_F16LE
	AV_CODEC_ID_PCM_F24LE          = C.AV_CODEC_ID_PCM_F24LE
	AV_CODEC_ID_PCM_VIDC           = C.AV_CODEC_ID_PCM_VIDC
	AV_CODEC_ID_ADPCM_IMA_QT       = C.AV_CODEC_ID_ADPCM_IMA_QT
	AV_CODEC_ID_ADPCM_IMA_WAV      = C.AV_CODEC_ID_ADPCM_IMA_WAV
	AV_CODEC_ID_ADPCM_IMA_DK3      = C.AV_CODEC_ID_ADPCM_IMA_DK3
	AV_CODEC_ID_ADPCM_IMA_DK4      = C.AV_CODEC_ID_ADPCM_IMA_DK4
	AV_CODEC_ID_ADPCM_IMA_WS       = C.AV_CODEC_ID_ADPCM_IMA_WS
	AV_CODEC_ID_ADPCM_IMA_SMJPEG   = C.AV_CODEC_ID_ADPCM_IMA_SMJPEG
	AV_CODEC_ID_ADPCM_MS           = C.AV_CODEC_ID_ADPCM_MS
	AV_CODEC_ID_ADPCM_4XM          = C.AV_CODEC_ID_ADPCM_4XM
	AV_CODEC_ID_ADPCM_XA           = C.AV_CODEC_ID_ADPCM_XA
	AV_CODEC_ID_ADPCM_ADX          = C.AV_CODEC_ID_ADPCM_ADX
	AV_CODEC_ID_ADPCM_EA           = C.AV_CODEC_ID_ADPCM_EA
	AV_CODEC_ID_ADPCM_G726         = C.AV_CODEC_ID_ADPCM_G726
	AV_CODEC_ID_ADPCM_CT           = C.AV_CODEC_ID_ADPCM_CT
	AV_CODEC_ID_ADPCM_SWF          = C.AV_CODEC_ID_ADPCM_SWF
	AV_CODEC_ID_ADPCM_YAMAHA       = C.AV_CODEC_ID_ADPCM_YAMAHA
	AV_CODEC_ID_ADPCM_SBPRO_4      = C.AV_CODEC_ID_ADPCM_SBPRO_4
	AV_CODEC_ID_ADPCM_SBPRO_3      = C.AV_CODEC_ID_ADPCM_SBPRO_3
	AV_CODEC_ID_ADPCM_SBPRO_2      = C.AV_CODEC_ID_ADPCM_SBPRO_2
	AV_CODEC_ID_ADPCM_THP          = C.AV_CODEC_ID_ADPCM_THP
	AV_CODEC_ID_ADPCM_IMA_AMV      = C.AV_CODEC_ID_ADPCM_IMA_AMV
	AV_CODEC_ID_ADPCM_EA_R1        = C.AV_CODEC_ID_ADPCM_EA_R1
	AV_CODEC_ID_ADPCM_EA_R3        = C.AV_CODEC_ID_ADPCM_EA_R3
	AV_CODEC_ID_ADPCM_EA_R2        = C.AV_CODEC_ID_ADPCM_EA_R2
	AV_CODEC_ID_ADPCM_IMA_EA_SEAD  = C.AV_CODEC_ID_ADPCM_IMA_EA_SEAD
	AV_CODEC_ID_ADPCM_IMA_EA_EACS  = C.AV_CODEC_ID_ADPCM_IMA_EA_EACS
	AV_CODEC_ID_ADPCM_EA_XAS       = C.AV_CODEC_ID_ADPCM_EA_XAS
	AV_CODEC_ID_ADPCM_EA_MAXIS_XA  = C.AV_CODEC_ID_ADPCM_EA_MAXIS_XA
	AV_CODEC_ID_ADPCM_IMA_ISS      = C.AV_CODEC_ID_ADPCM_IMA_ISS
	AV_CODEC_ID_ADPCM_G722         = C.AV_CODEC_ID_ADPCM_G722
	AV_CODEC_ID_ADPCM_IMA_APC      = C.AV_CODEC_ID_ADPCM_IMA_APC
	AV_CODEC_ID_ADPCM_VIMA         = C.AV_CODEC_ID_ADPCM_VIMA
	AV_CODEC_ID_ADPCM_AFC          = C.AV_CODEC_ID_ADPCM_AFC
	AV_CODEC_ID_ADPCM_IMA_OKI      = C.AV_CODEC_ID_ADPCM_IMA_OKI
	AV_CODEC_ID_ADPCM_DTK          = C.AV_CODEC_ID_ADPCM_DTK
	AV_CODEC_ID_ADPCM_IMA_RAD      = C.AV_CODEC_ID_ADPCM_IMA_RAD
	AV_CODEC_ID_ADPCM_G726LE       = C.AV_CODEC_ID_ADPCM_G726LE
	AV_CODEC_ID_ADPCM_THP_LE       = C.AV_CODEC_ID_ADPCM_THP_LE
	AV_CODEC_ID_ADPCM_PSX          = C.AV_CODEC_ID_ADPCM_PSX
	AV_CODEC_ID_ADPCM_AICA         = C.AV_CODEC_ID_ADPCM_AICA
	AV_CODEC_ID_ADPCM_IMA_DAT4     = C.AV_CODEC_ID_ADPCM_IMA_DAT4
	AV_CODEC_ID_ADPCM_MTAF         = C.AV_CODEC_ID_ADPCM_MTAF
	AV_CODEC_ID_AMR_NB             = C.AV_CODEC_ID_AMR_NB
	AV_CODEC_ID_AMR_WB             = C.AV_CODEC_ID_AMR_WB
	AV_CODEC_ID_RA_144             = C.AV_CODEC_ID_RA_144
	AV_CODEC_ID_RA_288             = C.AV_CODEC_ID_RA_288
	AV_CODEC_ID_ROQ_DPCM           = C.AV_CODEC_ID_ROQ_DPCM
	AV_CODEC_ID_INTERPLAY_DPCM     = C.AV_CODEC_ID_INTERPLAY_DPCM
	AV_CODEC_ID_XAN_DPCM           = C.AV_CODEC_ID_XAN_DPCM
	AV_CODEC_ID_SOL_DPCM           = C.AV_CODEC_ID_SOL_DPCM
	AV_CODEC_ID_SDX2_DPCM          = C.AV_CODEC_ID_SDX2_DPCM
	AV_CODEC_ID_GREMLIN_DPCM       = C.AV_CODEC_ID_GREMLIN_DPCM
	AV_CODEC_ID_MP2                = C.AV_CODEC_ID_MP2
	AV_CODEC_ID_MP3                = C.AV_CODEC_ID_MP3
	AV_CODEC_ID_AAC                = C.AV_CODEC_ID_AAC
	AV_CODEC_ID_AC3                = C.AV_CODEC_ID_AC3
	AV_CODEC_ID_DTS                = C.AV_CODEC_ID_DTS
	AV_CODEC_ID_VORBIS             = C.AV_CODEC_ID_VORBIS
	AV_CODEC_ID_DVAUDIO            = C.AV_CODEC_ID_DVAUDIO
	AV_CODEC_ID_WMAV1              = C.AV_CODEC_ID_WMAV1
	AV_CODEC_ID_WMAV2              = C.AV_CODEC_ID_WMAV2
	AV_CODEC_ID_MACE3              = C.AV_CODEC_ID_MACE3
	AV_CODEC_ID_MACE6              = C.AV_CODEC_ID_MACE6
	AV_CODEC_ID_VMDAUDIO           = C.AV_CODEC_ID_VMDAUDIO
	AV_CODEC_ID_FLAC               = C.AV_CODEC_ID_FLAC
	AV_CODEC_ID_MP3ADU             = C.AV_CODEC_ID_MP3ADU
	AV_CODEC_ID_MP3ON4             = C.AV_CODEC_ID_MP3ON4
	AV_CODEC_ID_SHORTEN            = C.AV_CODEC_ID_SHORTEN
	AV_CODEC_ID_ALAC               = C.AV_CODEC_ID_ALAC
	AV_CODEC_ID_WESTWOOD_SND1      = C.AV_CODEC_ID_WESTWOOD_SND1
	AV_CODEC_ID_GSM                = C.AV_CODEC_ID_GSM
	AV_CODEC_ID_QDM2               = C.AV_CODEC_ID_QDM2
	AV_CODEC_ID_COOK               = C.AV_CODEC_ID_COOK
	AV_CODEC_ID_TRUESPEECH         = C.AV_CODEC_ID_TRUESPEECH
	AV_CODEC_ID_TTA                = C.AV_CODEC_ID_TTA
	AV_CODEC_ID_SMACKAUDIO         = C.AV_CODEC_ID_SMACKAUDIO
	AV_CODEC_ID_QCELP              = C.AV_CODEC_ID_QCELP
	AV_CODEC_ID_WAVPACK            = C.AV_CODEC_ID_WAVPACK
	AV_CODEC_ID_DSICINAUDIO        = C.AV_CODEC_ID_DSICINAUDIO
	AV_CODEC_ID_IMC                = C.AV_CODEC_ID_IMC
	AV_CODEC_ID_MUSEPACK7          = C.AV_CODEC_ID_MUSEPACK7
	AV_CODEC_ID_MLP                = C.AV_CODEC_ID_MLP
	AV_CODEC_ID_GSM_MS             = C.AV_CODEC_ID_GSM_MS
	AV_CODEC_ID_ATRAC3             = C.AV_CODEC_ID_ATRAC3
	AV_CODEC_ID_APE                = C.AV_CODEC_ID_APE
	AV_CODEC_ID_NELLYMOSER         = C.AV_CODEC_ID_NELLYMOSER
	AV_CODEC_ID_MUSEPACK8          = C.AV_CODEC_ID_MUSEPACK8
	AV_CODEC_ID_SPEEX              = C.AV_CODEC_ID_SPEEX
	AV_CODEC_ID_WMAVOICE           = C.AV_CODEC_ID_WMAVOICE
	AV_CODEC_ID_WMAPRO             = C.AV_CODEC_ID_WMAPRO
	AV_CODEC_ID_WMALOSSLESS        = C.AV_CODEC_ID_WMALOSSLESS
	AV_CODEC_ID_ATRAC3P            = C.AV_CODEC_ID_ATRAC3P
	AV_CODEC_ID_EAC3               = C.AV_CODEC_ID_EAC3
	AV_CODEC_ID_SIPR               = C.AV_CODEC_ID_SIPR
	AV_CODEC_ID_MP1                = C.AV_CODEC_ID_MP1
	AV_CODEC_ID_TWINVQ             = C.AV_CODEC_ID_TWINVQ
	AV_CODEC_ID_TRUEHD             = C.AV_CODEC_ID_TRUEHD
	AV_CODEC_ID_MP4ALS             = C.AV_CODEC_ID_MP4ALS
	AV_CODEC_ID_ATRAC1             = C.AV_CODEC_ID_ATRAC1
	AV_CODEC_ID_BINKAUDIO_RDFT     = C.AV_CODEC_ID_BINKAUDIO_RDFT
	AV_CODEC_ID_BINKAUDIO_DCT      = C.AV_CODEC_ID_BINKAUDIO_DCT
	AV_CODEC_ID_AAC_LATM           = C.AV_CODEC_ID_AAC_LATM
	AV_CODEC_ID_QDMC               = C.AV_CODEC_ID_QDMC
	AV_CODEC_ID_CELT               = C.AV_CODEC_ID_CELT
	AV_CODEC_ID_G723_1             = C.AV_CODEC_ID_G723_1
	AV_CODEC_ID_G729               = C.AV_CODEC_ID_G729
	AV_CODEC_ID_8SVX_EXP           = C.AV_CODEC_ID_8SVX_EXP
	AV_CODEC_ID_8SVX_FIB           = C.AV_CODEC_ID_8SVX_FIB
	AV_CODEC_ID_BMV_AUDIO          = C.AV_CODEC_ID_BMV_AUDIO
	AV_CODEC_ID_RALF               = C.AV_CODEC_ID_RALF
	AV_CODEC_ID_IAC                = C.AV_CODEC_ID_IAC
	AV_CODEC_ID_ILBC               = C.AV_CODEC_ID_ILBC
	AV_CODEC_ID_OPUS               = C.AV_CODEC_ID_OPUS
	AV_CODEC_ID_COMFORT_NOISE      = C.AV_CODEC_ID_COMFORT_NOISE
	AV_CODEC_ID_TAK                = C.AV_CODEC_ID_TAK
	AV_CODEC_ID_METASOUND          = C.AV_CODEC_ID_METASOUND
	AV_CODEC_ID_PAF_AUDIO          = C.AV_CODEC_ID_PAF_AUDIO
	AV_CODEC_ID_ON2AVC             = C.AV_CODEC_ID_ON2AVC
	AV_CODEC_ID_DSS_SP             = C.AV_CODEC_ID_DSS_SP
	AV_CODEC_ID_CODEC2             = C.AV_CODEC_ID_CODEC2
	AV_CODEC_ID_FFWAVESYNTH        = C.AV_CODEC_ID_FFWAVESYNTH
	AV_CODEC_ID_SONIC              = C.AV_CODEC_ID_SONIC
	AV_CODEC_ID_SONIC_LS           = C.AV_CODEC_ID_SONIC_LS
	AV_CODEC_ID_EVRC               = C.AV_CODEC_ID_EVRC
	AV_CODEC_ID_SMV                = C.AV_CODEC_ID_SMV
	AV_CODEC_ID_DSD_LSBF           = C.AV_CODEC_ID_DSD_LSBF
	AV_CODEC_ID_DSD_MSBF           = C.AV_CODEC_ID_DSD_MSBF
	AV_CODEC_ID_DSD_LSBF_PLANAR    = C.AV_CODEC_ID_DSD_LSBF_PLANAR
	AV_CODEC_ID_DSD_MSBF_PLANAR    = C.AV_CODEC_ID_DSD_MSBF_PLANAR
	AV_CODEC_ID_4GV                = C.AV_CODEC_ID_4GV
	AV_CODEC_ID_INTERPLAY_ACM      = C.AV_CODEC_ID_INTERPLAY_ACM
	AV_CODEC_ID_XMA1               = C.AV_CODEC_ID_XMA1
	AV_CODEC_ID_XMA2               = C.AV_CODEC_ID_XMA2
	AV_CODEC_ID_DST                = C.AV_CODEC_ID_DST
	AV_CODEC_ID_ATRAC3AL           = C.AV_CODEC_ID_ATRAC3AL
	AV_CODEC_ID_ATRAC3PAL          = C.AV_CODEC_ID_ATRAC3PAL
	AV_CODEC_ID_DOLBY_E            = C.AV_CODEC_ID_DOLBY_E
	AV_CODEC_ID_APTX               = C.AV_CODEC_ID_APTX
	AV_CODEC_ID_APTX_HD            = C.AV_CODEC_ID_APTX_HD
	AV_CODEC_ID_SBC                = C.AV_CODEC_ID_SBC
	AV_CODEC_ID_ATRAC9             = C.AV_CODEC_ID_ATRAC9
	AV_CODEC_ID_FIRST_SUBTITLE     = C.AV_CODEC_ID_FIRST_SUBTITLE
	AV_CODEC_ID_DVD_SUBTITLE       = C.AV_CODEC_ID_DVD_SUBTITLE
	AV_CODEC_ID_DVB_SUBTITLE       = C.AV_CODEC_ID_DVB_SUBTITLE
	AV_CODEC_ID_TEXT               = C.AV_CODEC_ID_TEXT
	AV_CODEC_ID_XSUB               = C.AV_CODEC_ID_XSUB
	AV_CODEC_ID_SSA                = C.AV_CODEC_ID_SSA
	AV_CODEC_ID_MOV_TEXT           = C.AV_CODEC_ID_MOV_TEXT
	AV_CODEC_ID_HDMV_PGS_SUBTITLE  = C.AV_CODEC_ID_HDMV_PGS_SUBTITLE
	AV_CODEC_ID_DVB_TELETEXT       = C.AV_CODEC_ID_DVB_TELETEXT
	AV_CODEC_ID_SRT                = C.AV_CODEC_ID_SRT
	AV_CODEC_ID_MICRODVD           = C.AV_CODEC_ID_MICRODVD
	AV_CODEC_ID_EIA_608            = C.AV_CODEC_ID_EIA_608
	AV_CODEC_ID_JACOSUB            = C.AV_CODEC_ID_JACOSUB
	AV_CODEC_ID_SAMI               = C.AV_CODEC_ID_SAMI
	AV_CODEC_ID_REALTEXT           = C.AV_CODEC_ID_REALTEXT
	AV_CODEC_ID_STL                = C.AV_CODEC_ID_STL
	AV_CODEC_ID_SUBVIEWER1         = C.AV_CODEC_ID_SUBVIEWER1
	AV_CODEC_ID_SUBVIEWER          = C.AV_CODEC_ID_SUBVIEWER
	AV_CODEC_ID_SUBRIP             = C.AV_CODEC_ID_SUBRIP
	AV_CODEC_ID_WEBVTT             = C.AV_CODEC_ID_WEBVTT
	AV_CODEC_ID_MPL2               = C.AV_CODEC_ID_MPL2
	AV_CODEC_ID_VPLAYER            = C.AV_CODEC_ID_VPLAYER
	AV_CODEC_ID_PJS                = C.AV_CODEC_ID_PJS
	AV_CODEC_ID_ASS                = C.AV_CODEC_ID_ASS
	AV_CODEC_ID_HDMV_TEXT_SUBTITLE = C.AV_CODEC_ID_HDMV_TEXT_SUBTITLE
	AV_CODEC_ID_TTML               = C.AV_CODEC_ID_TTML
	AV_CODEC_ID_FIRST_UNKNOWN      = C.AV_CODEC_ID_FIRST_UNKNOWN
	AV_CODEC_ID_TTF                = C.AV_CODEC_ID_TTF
	AV_CODEC_ID_SCTE_35            = C.AV_CODEC_ID_SCTE_35
	AV_CODEC_ID_BINTEXT            = C.AV_CODEC_ID_BINTEXT
	AV_CODEC_ID_XBIN               = C.AV_CODEC_ID_XBIN
	AV_CODEC_ID_IDF                = C.AV_CODEC_ID_IDF
	AV_CODEC_ID_OTF                = C.AV_CODEC_ID_OTF
	AV_CODEC_ID_SMPTE_KLV          = C.AV_CODEC_ID_SMPTE_KLV
	AV_CODEC_ID_DVD_NAV            = C.AV_CODEC_ID_DVD_NAV
	AV_CODEC_ID_TIMED_ID3          = C.AV_CODEC_ID_TIMED_ID3
	AV_CODEC_ID_BIN_DATA           = C.AV_CODEC_ID_BIN_DATA
	AV_CODEC_ID_PROBE              = C.AV_CODEC_ID_PROBE
	AV_CODEC_ID_MPEG2TS            = C.AV_CODEC_ID_MPEG2TS
	AV_CODEC_ID_MPEG4SYSTEMS       = C.AV_CODEC_ID_MPEG4SYSTEMS
	AV_CODEC_ID_FFMETADATA         = C.AV_CODEC_ID_FFMETADATA
	AV_CODEC_ID_WRAPPED_AVFRAME    = C.AV_CODEC_ID_WRAPPED_AVFRAME
)

// SAMPLE FORMAT
const (
	AV_SAMPLE_FMT_NONE = C.AV_SAMPLE_FMT_NONE
	AV_SAMPLE_FMT_U8   = C.AV_SAMPLE_FMT_U8  ///< unsigned 8 bit
	AV_SAMPLE_FMT_S16  = C.AV_SAMPLE_FMT_S16 ///< signed 16 bit
	AV_SAMPLE_FMT_S32  = C.AV_SAMPLE_FMT_S32 ///< signed 32 bit
	AV_SAMPLE_FMT_FLT  = C.AV_SAMPLE_FMT_FLT ///< floa
	AV_SAMPLE_FMT_DBL  = C.AV_SAMPLE_FMT_DBL ///< doubl

	AV_SAMPLE_FMT_U8P  = C.AV_SAMPLE_FMT_U8P  ///< unsigned 8 bits, plana
	AV_SAMPLE_FMT_S16P = C.AV_SAMPLE_FMT_S16P ///< signed 16 bits, plana
	AV_SAMPLE_FMT_S32P = C.AV_SAMPLE_FMT_S32P ///< signed 32 bits, plana
	AV_SAMPLE_FMT_FLTP = C.AV_SAMPLE_FMT_FLTP ///< float, plana
	AV_SAMPLE_FMT_DBLP = C.AV_SAMPLE_FMT_DBLP ///< double, plana
	AV_SAMPLE_FMT_S64  = C.AV_SAMPLE_FMT_S64  ///< signed 64 bit
	AV_SAMPLE_FMT_S64P = C.AV_SAMPLE_FMT_S64P ///< signed 64 bits, plana

	AV_SAMPLE_FMT_NB = C.AV_SAMPLE_FMT_NB ///< Numb
)

const (
	AV_CH_LAYOUT_MONO              = C.AV_CH_LAYOUT_MONO
	AV_CH_LAYOUT_STEREO            = C.AV_CH_LAYOUT_STEREO
	AV_CH_LAYOUT_2POINT1           = C.AV_CH_LAYOUT_2POINT1
	AV_CH_LAYOUT_2_1               = C.AV_CH_LAYOUT_2_1
	AV_CH_LAYOUT_SURROUND          = C.AV_CH_LAYOUT_SURROUND
	AV_CH_LAYOUT_3POINT1           = C.AV_CH_LAYOUT_3POINT1
	AV_CH_LAYOUT_4POINT0           = C.AV_CH_LAYOUT_4POINT0
	AV_CH_LAYOUT_4POINT1           = C.AV_CH_LAYOUT_4POINT1
	AV_CH_LAYOUT_2_2               = C.AV_CH_LAYOUT_2_2
	AV_CH_LAYOUT_QUAD              = C.AV_CH_LAYOUT_QUAD
	AV_CH_LAYOUT_5POINT0           = C.AV_CH_LAYOUT_5POINT0
	AV_CH_LAYOUT_5POINT1           = C.AV_CH_LAYOUT_5POINT1
	AV_CH_LAYOUT_5POINT0_BACK      = C.AV_CH_LAYOUT_5POINT0_BACK
	AV_CH_LAYOUT_5POINT1_BACK      = C.AV_CH_LAYOUT_5POINT1_BACK
	AV_CH_LAYOUT_6POINT0           = C.AV_CH_LAYOUT_6POINT0
	AV_CH_LAYOUT_6POINT0_FRONT     = C.AV_CH_LAYOUT_6POINT0_FRONT
	AV_CH_LAYOUT_HEXAGONAL         = C.AV_CH_LAYOUT_HEXAGONAL
	AV_CH_LAYOUT_6POINT1           = C.AV_CH_LAYOUT_6POINT1
	AV_CH_LAYOUT_6POINT1_BACK      = C.AV_CH_LAYOUT_6POINT1_BACK
	AV_CH_LAYOUT_6POINT1_FRONT     = C.AV_CH_LAYOUT_6POINT1_FRONT
	AV_CH_LAYOUT_7POINT0           = C.AV_CH_LAYOUT_7POINT0
	AV_CH_LAYOUT_7POINT0_FRONT     = C.AV_CH_LAYOUT_7POINT0_FRONT
	AV_CH_LAYOUT_7POINT1           = C.AV_CH_LAYOUT_7POINT1
	AV_CH_LAYOUT_7POINT1_WIDE      = C.AV_CH_LAYOUT_7POINT1_WIDE
	AV_CH_LAYOUT_7POINT1_WIDE_BACK = C.AV_CH_LAYOUT_7POINT1_WIDE_BACK
	AV_CH_LAYOUT_OCTAGONAL         = C.AV_CH_LAYOUT_OCTAGONAL
	AV_CH_LAYOUT_HEXADECAGONAL     = C.AV_CH_LAYOUT_HEXADECAGONAL
	AV_CH_LAYOUT_STEREO_DOWNMIX    = C.AV_CH_LAYOUT_STEREO_DOWNMIX
)
