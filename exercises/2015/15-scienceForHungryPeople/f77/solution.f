C Solution for Advent of Code 2015 day 15.
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
C Distribute 100 teaspoons over the ingredients to maximise the cookie
C score: product of each property column total (negatives clamped to 0).
C Part 2 additionally requires the calorie total to equal exactly 500.
C This is a dense integer optimisation search — Fortran's home turf.

C     Shared search state. CAP/DUR/FLA/TEX/CAL are per-ingredient property
C     columns; N the ingredient count; GOAL the calorie target (-1 = none);
C     AMT the current recipe; BEST the running maximum score.
      BLOCK DATA DAY15INIT
      INTEGER N, GOAL, CAP, DUR, FLA, TEX, CAL, AMT, BEST
      COMMON /D15/ CAP(20), DUR(20), FLA(20), TEX(20), CAL(20),
     &             AMT(20), N, GOAL, BEST
      DATA N /0/, GOAL /-1/, BEST /0/
      END

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER N, GOAL, CAP, DUR, FLA, TEX, CAL, AMT, BEST
      COMMON /D15/ CAP(20), DUR(20), FLA(20), TEX(20), CAL(20),
     &             AMT(20), N, GOAL, BEST

      CALL LOAD(INPUT_PATH)
      GOAL = -1
      BEST = 0
      CALL REC(1, 100)
      CALL ITOA(BEST, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER N, GOAL, CAP, DUR, FLA, TEX, CAL, AMT, BEST
      COMMON /D15/ CAP(20), DUR(20), FLA(20), TEX(20), CAL(20),
     &             AMT(20), N, GOAL, BEST

      CALL LOAD(INPUT_PATH)
      GOAL = 500
      BEST = 0
      CALL REC(1, 100)
      CALL ITOA(BEST, ANSWER, ANSWER_LEN)
      END

C     REC — enumerate compositions of REMAIN teaspoons across ingredients
C     I..N. The last ingredient takes whatever is left; each complete recipe
C     is scored and BEST updated. Recursive (gfortran extension over F77).
      RECURSIVE SUBROUTINE REC(I, REMAIN)
      INTEGER I, REMAIN
      INTEGER N, GOAL, CAP, DUR, FLA, TEX, CAL, AMT, BEST
      COMMON /D15/ CAP(20), DUR(20), FLA(20), TEX(20), CAL(20),
     &             AMT(20), N, GOAL, BEST
      INTEGER A, S, SCORE

      IF (I .EQ. N) THEN
        AMT(I) = REMAIN
        S = SCORE()
        IF (S .GT. BEST) BEST = S
        RETURN
      END IF
      DO 10 A = 0, REMAIN
        AMT(I) = A
        CALL REC(I + 1, REMAIN - A)
10    CONTINUE
      END

C     SCORE — product of the four property column totals (negatives clamped
C     to 0). Returns 0 if a calorie goal is set and unmet.
      INTEGER FUNCTION SCORE()
      INTEGER N, GOAL, CAP, DUR, FLA, TEX, CAL, AMT, BEST
      COMMON /D15/ CAP(20), DUR(20), FLA(20), TEX(20), CAL(20),
     &             AMT(20), N, GOAL, BEST
      INTEGER TC, TD, TF, TT, TK, I

      TC = 0
      TD = 0
      TF = 0
      TT = 0
      TK = 0
      DO 10 I = 1, N
        TC = TC + CAP(I) * AMT(I)
        TD = TD + DUR(I) * AMT(I)
        TF = TF + FLA(I) * AMT(I)
        TT = TT + TEX(I) * AMT(I)
        TK = TK + CAL(I) * AMT(I)
10    CONTINUE
      IF (GOAL .GE. 0 .AND. TK .NE. GOAL) THEN
        SCORE = 0
        RETURN
      END IF
      IF (TC .LT. 0) TC = 0
      IF (TD .LT. 0) TD = 0
      IF (TF .LT. 0) TF = 0
      IF (TT .LT. 0) TT = 0
      SCORE = TC * TD * TF * TT
      END

C     LOAD — parse each ingredient's five signed integers. Non-digits are
C     blanked, except a '-' that immediately precedes a digit (so negative
C     property values survive), then list-directed read pulls the values.
      SUBROUTINE LOAD(INPUT_PATH)
      CHARACTER*(*) INPUT_PATH
      INTEGER N, GOAL, CAP, DUR, FLA, TEX, CAL, AMT, BEST
      COMMON /D15/ CAP(20), DUR(20), FLA(20), TEX(20), CAL(20),
     &             AMT(20), N, GOAL, BEST
      CHARACTER*256 LINE
      CHARACTER C1, C2
      INTEGER IOS, I

      N = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 100
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 20
      DO 30 I = 1, 256
        C1 = LINE(I:I)
        IF (C1 .GE. '0' .AND. C1 .LE. '9') GOTO 30
        IF (C1 .EQ. '-' .AND. I .LT. 256) THEN
          C2 = LINE(I+1:I+1)
          IF (C2 .GE. '0' .AND. C2 .LE. '9') GOTO 30
        END IF
        LINE(I:I) = ' '
30    CONTINUE
      N = N + 1
      READ(LINE, *) CAP(N), DUR(N), FLA(N), TEX(N), CAL(N)
      GOTO 20
100   CLOSE(10)
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
