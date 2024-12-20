#include <libipset/ipset.h>
#include <stdio.h>
#include <stdarg.h>
#include "cipset.h"

void goips_init(void) {
	ipset_load_types();
}

static int goips_custom_errorfn(struct ipset *ipset, void *p, int status, const char *msg, ...) {
	char buffer[8192];
	va_list args;
	va_start(args, msg);
	vsnprintf(buffer, 8192 ,msg, args);
	va_end (args);

	fprintf(stderr, "+++1 goips_custom_errorfn: %p %d", p, status);
	switch (status) {
	case IPSET_NO_PROBLEM:
		fprintf(stderr, " %s", "no problem"); break;
	case IPSET_OTHER_PROBLEM:
		fprintf(stderr, " %s", "other problem"); break;
	case IPSET_PARAMETER_PROBLEM:
		fprintf(stderr, " %s", "parameter problem"); break;
	case IPSET_VERSION_PROBLEM:
		fprintf(stderr, " %s", "version problem"); break;
	case IPSET_SESSION_PROBLEM:
		fprintf(stderr, " %s", "session problem"); break;
	}

	fprintf(stderr, " msg: %s\n", buffer);
	fflush(stderr);

	return status;
}

static int goips_standard_errorfn(struct ipset *ipset, void *p) {
	fprintf(stderr, "+++2 goips_standard_errorfn: %p\n", p);

	struct ipset_session *session = ipset_session(ipset);
	bool is_interactive = ipset_is_interactive(ipset);
	enum ipset_err_type err_type = ipset_session_report_type(session);

	fprintf(stderr, "  is interactive: %d\n", is_interactive);
	fprintf(stderr, "  err_type: %d ", err_type);
	switch (err_type) {
	case IPSET_NO_ERROR:
		fprintf(stderr, "no error\n"); break;
	case IPSET_NOTICE:
		fprintf(stderr, "notice\n"); break;
	case IPSET_WARNING:
		fprintf(stderr, "warning\n"); break;
	case IPSET_ERROR:
		fprintf(stderr, "error\n"); break;
	}

	fprintf(stderr, "  msg: %s", ipset_session_report_msg(session));

	ipset_session_report_reset(session);
	return -1;
}

static int goips_print_outfn(struct ipset_session *session, void *p, const char *fmt, ...) {
	fprintf(stderr, "+++3 goips_print_outfn: %p\n", p);

	char buffer[8192];
	va_list args;
	va_start(args, fmt);
	vsnprintf(buffer, 8192 , fmt, args);
	va_end (args);

	fprintf(stderr, "%s\n", buffer);

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
