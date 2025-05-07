INCLUDE "targets/atsamd21.ld"

/* Use fixed starting point for user data (highest 32kB) */
__flash_data_start = 0x38000;
