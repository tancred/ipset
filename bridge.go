package ipset

/*
#include <stdarg.h>
#include <libipset/ipset.h>

extern void goipsStandardErrorFn(struct ipset *ipset, void *p, int errType, const char *msg);
extern int goipsCustomErrorFn(struct ipset *ipset, void *p, int status, const char *msg);
extern void goipsPrintOutFn(void *p, const char *msg);

int goips_custom_errorfn(struct ipset *ipset, void *p, int status, const char *msg, ...) {
	char buffer[8192];
	va_list args;
	va_start(args, msg);
	vsnprintf(buffer, 8192 ,msg, args);
	va_end (args);
	return goipsCustomErrorFn(ipset, p, status, buffer);
}

int goips_standard_errorfn(struct ipset *ipset, void *p) {
	struct ipset_session *session = ipset_session(ipset);
	enum ipset_err_type err_type = ipset_session_report_type(session);

	const char *msg = ipset_session_report_msg(session);
	goipsStandardErrorFn(ipset, p, err_type, msg);
	ipset_session_report_reset(session);

	return -1;
}

int goips_print_outfn(struct ipset_session *session, void *p, const char *fmt, ...) {
	char buffer[8192];
	va_list args;
	va_start(args, fmt);
	vsnprintf(buffer, 8192 , fmt, args);
	va_end (args);

	goipsPrintOutFn(p, buffer);
	return 0;
}

int goips_custom_printf(struct ipset *ipset, void *p) {
	return ipset_custom_printf(
		ipset,
		goips_custom_errorfn,
		goips_standard_errorfn,
		goips_print_outfn,
		p
	);
}

*/
import "C"
