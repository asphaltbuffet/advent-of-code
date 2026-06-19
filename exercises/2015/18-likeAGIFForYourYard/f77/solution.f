C Solution for Advent of Code 2015 day 18.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Conway's Game of Life on a light grid. Part 2 keeps the four corners
C stuck on. A native 2D integer array with the grid's exact 0-based bounds
C is ping-ponged each generation — the dense-grid case Fortran handles best.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER SOLVE
      CALL ITOA(SOLVE(INPUT_PATH, 0), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER SOLVE
      CALL ITOA(SOLVE(INPUT_PATH, 1), ANSWER, ANSWER_LEN)
      END

C     SOLVE — animate the grid and return the final lit count. STUCK forces
C     the four corners on (initially and after each step). Steps: the small
C     example runs 4 (part 1) or 5 (part 2); the real grid runs 100.
      INTEGER FUNCTION SOLVE(INPUT_PATH, STUCK)
      CHARACTER*(*) INPUT_PATH
      INTEGER STUCK
      INTEGER G(0:99,0:99), NX(0:99,0:99)
      INTEGER R, C, NR, NC, IOS, I, J, DR, DC, ON, S, STEPS, TOT
      CHARACTER*200 LINE
      SAVE G, NX

      R = 0
      C = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 30
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 20
      C = LEN_TRIM(LINE)
      DO 25 J = 1, C
        IF (LINE(J:J) .EQ. '#') THEN
          G(R, J-1) = 1
        ELSE
          G(R, J-1) = 0
        END IF
25    CONTINUE
      R = R + 1
      GOTO 20
30    CLOSE(10)

      IF (R .GT. 6) THEN
        STEPS = 100
      ELSE IF (STUCK .EQ. 1) THEN
        STEPS = 5
      ELSE
        STEPS = 4
      END IF

      IF (STUCK .EQ. 1) CALL STICK(G, R, C)

      DO 80 S = 1, STEPS
        DO 60 I = 0, R - 1
          DO 55 J = 0, C - 1
            ON = 0
            DO 50 DR = -1, 1
              DO 45 DC = -1, 1
                IF (DR .EQ. 0 .AND. DC .EQ. 0) GOTO 45
                NR = I + DR
                NC = J + DC
                IF (NR .GE. 0 .AND. NR .LT. R .AND.
     &              NC .GE. 0 .AND. NC .LT. C) THEN
                  IF (G(NR, NC) .EQ. 1) ON = ON + 1
                END IF
45            CONTINUE
50          CONTINUE
            IF (G(I, J) .EQ. 1) THEN
              IF (ON .EQ. 2 .OR. ON .EQ. 3) THEN
                NX(I, J) = 1
              ELSE
                NX(I, J) = 0
              END IF
            ELSE
              IF (ON .EQ. 3) THEN
                NX(I, J) = 1
              ELSE
                NX(I, J) = 0
              END IF
            END IF
55        CONTINUE
60      CONTINUE
        DO 70 I = 0, R - 1
          DO 65 J = 0, C - 1
            G(I, J) = NX(I, J)
65        CONTINUE
70      CONTINUE
        IF (STUCK .EQ. 1) CALL STICK(G, R, C)
80    CONTINUE

      TOT = 0
      DO 95 I = 0, R - 1
        DO 90 J = 0, C - 1
          TOT = TOT + G(I, J)
90      CONTINUE
95    CONTINUE
      SOLVE = TOT
      END

C     STICK — force the four corner lights on.
      SUBROUTINE STICK(G, R, C)
      INTEGER G(0:99,0:99), R, C
      G(0, 0) = 1
      G(0, C-1) = 1
      G(R-1, 0) = 1
      G(R-1, C-1) = 1
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
