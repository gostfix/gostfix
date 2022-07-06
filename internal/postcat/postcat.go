package postcat

/* Application-specific. */
const (
	PC_FLAG_SEARCH_QUEUE    = (1 << 0) /* search queue */
	PC_FLAG_PRINT_OFFSET    = (1 << 1) /* print record offsets */
	PC_FLAG_PRINT_ENV       = (1 << 2) /* print envelope records */
	PC_FLAG_PRINT_HEADER    = (1 << 3) /* print header records */
	PC_FLAG_PRINT_BODY      = (1 << 4) /* print body records */
	PC_FLAG_PRINT_RTYPE_DEC = (1 << 5) /* print decimal record type */
	PC_FLAG_PRINT_RTYPE_SYM = (1 << 6) /* print symbolic record type */
	PC_FLAG_RAW             = (1 << 7) /* don't follow pointers */
)

const (
	PC_MASK_PRINT_TEXT = (PC_FLAG_PRINT_HEADER | PC_FLAG_PRINT_BODY)
	PC_MASK_PRINT_ALL  = (PC_FLAG_PRINT_ENV | PC_MASK_PRINT_TEXT)
)

/*
 * State machine.
 */
const (
	PC_STATE_ENV    = 0 /* initial or extracted envelope */
	PC_STATE_HEADER = 1 /* primary header */
	PC_STATE_BODY   = 2 /* other */
)
