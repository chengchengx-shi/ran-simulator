/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "E2AP-IEs"
 * 	found in "e2ap-v01.00.00.asn"
 * 	`asn1c -fcompound-names -fincludes-quoted -fno-include-deps -findirect-choice -gen-PER -no-gen-OER -D.`
 */

#ifndef	_RICactionType_H_
#define	_RICactionType_H_


#include "asn_application.h"

/* Including external dependencies */
#include "NativeEnumerated.h"

#ifdef __cplusplus
extern "C" {
#endif

/* Dependencies */
typedef enum RICactionType {
	RICactionType_report	= 0,
	RICactionType_insert	= 1,
	RICactionType_policy	= 2
	/*
	 * Enumeration is extensible
	 */
} e_RICactionType;

/* RICactionType */
typedef long	 RICactionType_t;

/* Implementation */
extern asn_per_constraints_t asn_PER_type_RICactionType_constr_1;
extern asn_TYPE_descriptor_t asn_DEF_RICactionType;
extern const asn_INTEGER_specifics_t asn_SPC_RICactionType_specs_1;
asn_struct_free_f RICactionType_free;
asn_struct_print_f RICactionType_print;
asn_constr_check_f RICactionType_constraint;
ber_type_decoder_f RICactionType_decode_ber;
der_type_encoder_f RICactionType_encode_der;
xer_type_decoder_f RICactionType_decode_xer;
xer_type_encoder_f RICactionType_encode_xer;
per_type_decoder_f RICactionType_decode_uper;
per_type_encoder_f RICactionType_encode_uper;
per_type_decoder_f RICactionType_decode_aper;
per_type_encoder_f RICactionType_encode_aper;

#ifdef __cplusplus
}
#endif

#endif	/* _RICactionType_H_ */
#include "asn_internal.h"
