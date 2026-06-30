C Solution for Advent of Code 2015 day 1.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C elf's generated C harness calls these subroutines, handles the wire
C protocol, timing, and temp-file management — you only edit this file.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 are code.
C Column 1 'C' is a comment. Column 6 non-space is a continuation line.
C
C The input is parentheses: '(' goes up a floor, ')' goes down. Part 1 is
C the final floor; part 2 is the 1-based position of the first char that
C first sends Santa to floor -1 (the basement).

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER FLOOR, POS

      CALL SOLVE(INPUT_PATH, FLOOR, POS)
      CALL ITOA(FLOOR, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER FLOOR, POS

      CALL SOLVE(INPUT_PATH, FLOOR, POS)
      CALL ITOA(POS, ANSWER, ANSWER_LEN)
      END

C     SOLVE — read the input one character at a time, tracking the running
C     floor. Returns the final FLOOR and POS, the 1-based index of the first
C     character to reach floor -1 (0 if the basement is never entered).
      SUBROUTINE SOLVE(INPUT_PATH, FLOOR, POS)
      CHARACTER*(*) INPUT_PATH
      INTEGER FLOOR, POS
      CHARACTER CH
      INTEGER I, IOS

      FLOOR = 0
      POS = 0
      I = 0

C     Non-advancing, one-char-at-a-time read avoids buffering the whole
C     7000-char single line. IOSTAT traps end-of-record and end-of-file.
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD', ACCESS='SEQUENTIAL')
10    READ(10, '(A1)', ADVANCE='NO', IOSTAT=IOS) CH
      IF (IOS .NE. 0) GOTO 20
      IF (CH .EQ. '(') THEN
         FLOOR = FLOOR + 1
         I = I + 1
      ELSE IF (CH .EQ. ')') THEN
         FLOOR = FLOOR - 1
         I = I + 1
      ELSE
C        Ignore any stray whitespace/newline characters.
         GOTO 10
      END IF
      IF (FLOOR .EQ. -1 .AND. POS .EQ. 0) POS = I
      GOTO 10
20    CLOSE(10)
      END

C     ITOA — format integer N into S (left-justified) and report its length.
      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF

      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
