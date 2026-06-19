C Solution for Advent of Code 2015 day 6.
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
C A 1000x1000 light grid with rectangular range commands. Part 1 treats
C lights as on/off (count lit); part 2 as brightness (turn on +1, turn
C off -1 clamped at 0, toggle +2; sum brightness). This is the kind of
C dense numeric grid work Fortran's contiguous arrays handle natively.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER TOTAL

      CALL SOLVE(INPUT_PATH, 1, TOTAL)
      CALL ITOA(TOTAL, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER TOTAL

      CALL SOLVE(INPUT_PATH, 2, TOTAL)
      CALL ITOA(TOTAL, ANSWER, ANSWER_LEN)
      END

C     SOLVE — apply every instruction to the grid, then sum it. MODE 1 is the
C     on/off model, MODE 2 the brightness model. Returns the grand total.
      SUBROUTINE SOLVE(INPUT_PATH, MODE, TOTAL)
      CHARACTER*(*) INPUT_PATH
      INTEGER MODE, TOTAL
      INTEGER G(0:999,0:999)
      CHARACTER*256 LINE
      INTEGER X1, Y1, X2, Y2, X, Y, ACT, IOS, I

C     ACT: 1 = turn on, 0 = turn off, 2 = toggle.
      DO 10 Y = 0, 999
        DO 5 X = 0, 999
          G(X,Y) = 0
5       CONTINUE
10    CONTINUE

      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 100

C     Identify the verb from its fixed-length prefix.
      IF (LINE(1:7) .EQ. 'turn on') THEN
        ACT = 1
      ELSE IF (LINE(1:8) .EQ. 'turn off') THEN
        ACT = 0
      ELSE IF (LINE(1:6) .EQ. 'toggle') THEN
        ACT = 2
      ELSE
        GOTO 20
      END IF

C     Strip everything but digits to spaces, then read the four coordinates
C     with list-directed input regardless of the verb's length.
      DO 30 I = 1, 256
        IF (LINE(I:I) .LT. '0' .OR. LINE(I:I) .GT. '9')
     &      LINE(I:I) = ' '
30    CONTINUE
      READ(LINE, *) X1, Y1, X2, Y2

      DO 50 Y = Y1, Y2
        DO 40 X = X1, X2
          IF (MODE .EQ. 1) THEN
            IF (ACT .EQ. 1) THEN
              G(X,Y) = 1
            ELSE IF (ACT .EQ. 0) THEN
              G(X,Y) = 0
            ELSE
              G(X,Y) = 1 - G(X,Y)
            END IF
          ELSE
            IF (ACT .EQ. 1) THEN
              G(X,Y) = G(X,Y) + 1
            ELSE IF (ACT .EQ. 0) THEN
              IF (G(X,Y) .GT. 0) G(X,Y) = G(X,Y) - 1
            ELSE
              G(X,Y) = G(X,Y) + 2
            END IF
          END IF
40      CONTINUE
50    CONTINUE
      GOTO 20

100   CLOSE(10)

      TOTAL = 0
      DO 120 Y = 0, 999
        DO 110 X = 0, 999
          TOTAL = TOTAL + G(X,Y)
110     CONTINUE
120   CONTINUE
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
