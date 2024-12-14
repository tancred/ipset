#include <libipset/ipset.h>
#include "cipset.h"

void goips_init(void) {
	ipset_load_types();
}
