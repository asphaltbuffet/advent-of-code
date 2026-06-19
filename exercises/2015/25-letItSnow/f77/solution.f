C Solution for Advent of Code 2015 day 25.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Codes fill the grid diagonally; cell (row, col) is the n-th code, where
C n = T(row+col-2) + col. The n-th code is 20151125 * 252533^(n-1) mod
C 33554393. The triangular index and modexp intermediates need INTEGER*8.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER ROW, COL
      INTEGER*8 DIAG, N, RESULT, BASE, EXPN
      INTEGER*8 MD, FIRST, MULT

      MD = 33554393_8
      FIRST = 20151125_8
      MULT = 252533_8

      CALL READRC(INPUT_PATH, ROW, COL)
      DIAG = INT(ROW + COL - 2, 8)
      N = DIAG * (DIAG + 1) / 2 + INT(COL, 8)

C     RESULT = MULT^(N-1) mod MD by fast exponentiation.
      RESULT = 1
      BASE = MOD(MULT, MD)
      EXPN = N - 1
10    IF (EXPN .GT. 0) THEN
        IF (IAND(EXPN, 1_8) .EQ. 1) RESULT = MOD(RESULT * BASE, MD)
        BASE = MOD(BASE * BASE, MD)
        EXPN = ISHFT(EXPN, -1)
        GOTO 10
      END IF

      CALL ITOA8(MOD(FIRST * RESULT, MD), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
C     Day 25 has no part 2 — the final star comes from finishing the rest.
      ANSWER = 'Merry Christmas!'
      ANSWER_LEN = 16
      END

C     READRC — parse the row and column integers from the puzzle prose.
      SUBROUTINE READRC(INPUT_PATH, ROW, COL)
      CHARACTER*(*) INPUT_PATH
      INTEGER ROW, COL
      CHARACTER*256 LINE
      INTEGER I
      CHARACTER C

      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
      READ(10, '(A)') LINE
      CLOSE(10)
      DO 10 I = 1, 256
        C = LINE(I:I)
        IF (C .LT. '0' .OR. C .GT. '9') LINE(I:I) = ' '
10    CONTINUE
      READ(LINE, *) ROW, COL
      END

      SUBROUTINE ITOA8(N, S, L)
      INTEGER*8 N
      INTEGER L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
