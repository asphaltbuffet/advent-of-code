C Solution for Advent of Code 2015 day 3.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C A string of ^v<> moves. Part 1: distinct houses Santa visits. Part 2:
C Santa and Robo-Santa take alternate moves; distinct houses either visits.
C No hash set in F77, so visited houses live in a fixed 2D grid keyed by
C (x,y) offset; AoC walks stay well within +/-1024 of the origin.

      INTEGER FUNCTION COUNTH(INPUT_PATH, NSANTA)
      CHARACTER*(*) INPUT_PATH
      INTEGER NSANTA
      INTEGER VIS(-1024:1024,-1024:1024)
      INTEGER X(2), Y(2), WHO, CNT, DX, DY, IOS, I, N
      CHARACTER*16384 LINE
      CHARACTER CH
C     Keep the 16 MB visited grid in static storage, not on the stack —
C     the runner subprocess may have a small stack ulimit.
      SAVE VIS

      DO 20 DY = -1024, 1024
        DO 10 DX = -1024, 1024
          VIS(DX,DY) = 0
10      CONTINUE
20    CONTINUE

      X(1) = 0
      Y(1) = 0
      X(2) = 0
      Y(2) = 0
      WHO = 1
      CNT = 1
      VIS(0,0) = 1

C     The move string is one long line; read it whole, then scan characters.
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
      READ(10, '(A)', IOSTAT=IOS) LINE
      CLOSE(10)
      IF (IOS .NE. 0) GOTO 100
      N = LEN_TRIM(LINE)

      DO 50 I = 1, N
        CH = LINE(I:I)
        DX = 0
        DY = 0
        IF (CH .EQ. '^') THEN
          DY = 1
        ELSE IF (CH .EQ. 'v') THEN
          DY = -1
        ELSE IF (CH .EQ. '>') THEN
          DX = 1
        ELSE IF (CH .EQ. '<') THEN
          DX = -1
        ELSE
          GOTO 50
        END IF
        X(WHO) = X(WHO) + DX
        Y(WHO) = Y(WHO) + DY
        IF (VIS(X(WHO),Y(WHO)) .EQ. 0) THEN
          VIS(X(WHO),Y(WHO)) = 1
          CNT = CNT + 1
        END IF
        IF (NSANTA .EQ. 2) WHO = 3 - WHO
50    CONTINUE

100   COUNTH = CNT
      END

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER COUNTH
      CALL ITOA(COUNTH(INPUT_PATH, 1), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER COUNTH
      CALL ITOA(COUNTH(INPUT_PATH, 2), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
