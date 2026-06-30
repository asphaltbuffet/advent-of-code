C Solution for Advent of Code 2015 day 17.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Count subsets of containers that hold exactly TARGET liters. Part 1 is
C the total; part 2 is how many use the fewest containers. A 2^n bitmask
C sweep tallies a histogram of valid subsets by container count.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER CNT(0:20), N, TOTAL, K

      CALL SOLVE(INPUT_PATH, CNT, N)
      TOTAL = 0
      DO 10 K = 0, N
        TOTAL = TOTAL + CNT(K)
10    CONTINUE
      CALL ITOA(TOTAL, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER CNT(0:20), N, K, ANS

      CALL SOLVE(INPUT_PATH, CNT, N)
C     First non-zero bucket is the minimum-container count.
      ANS = 0
      DO 10 K = 1, N
        IF (CNT(K) .GT. 0) THEN
          ANS = CNT(K)
          GOTO 20
        END IF
10    CONTINUE
20    CALL ITOA(ANS, ANSWER, ANSWER_LEN)
      END

C     SOLVE — sweep every subset of the containers, building CNT(k) = number
C     of k-container subsets that sum to exactly the target volume.
      SUBROUTINE SOLVE(INPUT_PATH, CNT, N)
      CHARACTER*(*) INPUT_PATH
      INTEGER CNT(0:20), N
      INTEGER SZ(20), IOS, V, MASK, I, SUM, BITS, TARGET, NMASK

      N = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, *, IOSTAT=IOS) V
      IF (IOS .NE. 0) GOTO 30
      N = N + 1
      SZ(N) = V
      GOTO 20
30    CLOSE(10)

      IF (N .LE. 5) THEN
        TARGET = 25
      ELSE
        TARGET = 150
      END IF

      DO 40 I = 0, N
        CNT(I) = 0
40    CONTINUE

      NMASK = ISHFT(1, N)
      DO 60 MASK = 0, NMASK - 1
        SUM = 0
        BITS = 0
        DO 50 I = 0, N - 1
          IF (BTEST(MASK, I)) THEN
            SUM = SUM + SZ(I + 1)
            BITS = BITS + 1
          END IF
50      CONTINUE
        IF (SUM .EQ. TARGET) CNT(BITS) = CNT(BITS) + 1
60    CONTINUE
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
