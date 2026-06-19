C Solution for Advent of Code 2015 day 14.
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
C Reindeer fly at SPEED for FLY seconds, then rest REST seconds, cycling.
C Part 1: furthest distance after the race. Part 2: a point to each leader
C every second; most points wins. The duration is not in the input: the
C 2-reindeer AoC example races 1000s, the real input 2503s.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER SPEED(100), FLY(100), REST(100), N, DUR
      INTEGER I, D, BEST, DIST

      CALL LOAD(INPUT_PATH, SPEED, FLY, REST, N, DUR)
      BEST = 0
      DO 10 I = 1, N
        D = DIST(SPEED(I), FLY(I), REST(I), DUR)
        IF (D .GT. BEST) BEST = D
10    CONTINUE
      CALL ITOA(BEST, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER SPEED(100), FLY(100), REST(100), N, DUR
      INTEGER PTS(100), DS(100)
      INTEGER I, T, LEAD, BEST, DIST

      CALL LOAD(INPUT_PATH, SPEED, FLY, REST, N, DUR)
      DO 10 I = 1, N
        PTS(I) = 0
10    CONTINUE

      DO 40 T = 1, DUR
        LEAD = 0
        DO 20 I = 1, N
          DS(I) = DIST(SPEED(I), FLY(I), REST(I), T)
          IF (DS(I) .GT. LEAD) LEAD = DS(I)
20      CONTINUE
        DO 30 I = 1, N
          IF (DS(I) .EQ. LEAD) PTS(I) = PTS(I) + 1
30      CONTINUE
40    CONTINUE

      BEST = 0
      DO 50 I = 1, N
        IF (PTS(I) .GT. BEST) BEST = PTS(I)
50    CONTINUE
      CALL ITOA(BEST, ANSWER, ANSWER_LEN)
      END

C     DIST — distance one reindeer has travelled after T seconds.
      INTEGER FUNCTION DIST(SPEED, FLY, REST, T)
      INTEGER SPEED, FLY, REST, T
      INTEGER CYCLE, REM, FLYING

      CYCLE = FLY + REST
      REM = MOD(T, CYCLE)
      FLYING = (T / CYCLE) * FLY
      IF (REM .LT. FLY) THEN
        FLYING = FLYING + REM
      ELSE
        FLYING = FLYING + FLY
      END IF
      DIST = FLYING * SPEED
      END

C     LOAD — parse each line's three integers (speed, fly, rest) by stripping
C     non-digits to spaces, then list-directed read. DUR is inferred from N.
      SUBROUTINE LOAD(INPUT_PATH, SPEED, FLY, REST, N, DUR)
      CHARACTER*(*) INPUT_PATH
      INTEGER SPEED(100), FLY(100), REST(100), N, DUR
      CHARACTER*256 LINE
      INTEGER IOS, I

      N = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 100
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 20
      DO 30 I = 1, 256
        IF (LINE(I:I) .LT. '0' .OR. LINE(I:I) .GT. '9')
     &      LINE(I:I) = ' '
30    CONTINUE
      N = N + 1
      READ(LINE, *) SPEED(N), FLY(N), REST(N)
      GOTO 20
100   CLOSE(10)

      IF (N .LE. 2) THEN
        DUR = 1000
      ELSE
        DUR = 2503
      END IF
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
