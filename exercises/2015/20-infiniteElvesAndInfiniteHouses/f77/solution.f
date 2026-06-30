C Solution for Advent of Code 2015 day 20.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Elf n delivers presents to multiples of n. Part 1: 10*n to every multiple
C (house h gets 10 * sum-of-divisors). Part 2: 11*n to only the first 50.
C A divisor sieve over target/10 houses — dense integer array work, exactly
C Fortran's strength.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER SIEVE
      CALL ITOA(SIEVE(INPUT_PATH, 1), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER SIEVE
      CALL ITOA(SIEVE(INPUT_PATH, 2), ANSWER, ANSWER_LEN)
      END

C     SIEVE — accumulate each elf's gifts into the houses array, then return
C     the lowest house meeting the target. PART 1 visits all multiples (10*n);
C     PART 2 visits only the first 50 (11*n).
      INTEGER FUNCTION SIEVE(INPUT_PATH, PART)
      CHARACTER*(*) INPUT_PATH
      INTEGER PART
      INTEGER MAXH
      PARAMETER (MAXH = 4000000)
      INTEGER H(MAXH)
      INTEGER TARGET, LIMIT, N, I, C, IOS
      SAVE H

      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
      READ(10, *, IOSTAT=IOS) TARGET
      CLOSE(10)

      LIMIT = TARGET / 10 + 1
      IF (LIMIT .GT. MAXH) LIMIT = MAXH

      DO 10 I = 1, LIMIT
        H(I) = 0
10    CONTINUE

      IF (PART .EQ. 1) THEN
        DO 30 N = 1, LIMIT
          DO 20 I = N, LIMIT, N
            H(I) = H(I) + 10 * N
20        CONTINUE
30      CONTINUE
      ELSE
        DO 50 N = 1, LIMIT
          C = 0
          I = N
40        IF (I .LE. LIMIT .AND. C .LT. 50) THEN
            H(I) = H(I) + 11 * N
            I = I + N
            C = C + 1
            GOTO 40
          END IF
50      CONTINUE
      END IF

      DO 60 I = 1, LIMIT
        IF (H(I) .GE. TARGET) THEN
          SIEVE = I
          RETURN
        END IF
60    CONTINUE
      SIEVE = -1
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
