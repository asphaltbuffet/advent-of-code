C Solution for Advent of Code 2015 day 10.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Look-and-say: repeatedly replace runs of a digit with "<count><digit>".
C Part 1 runs 40 iterations, part 2 runs 50, reporting the result length.
C The string grows ~30%/step to several million chars, so two large static
C buffers are ping-ponged (lengths tracked explicitly, never LEN_TRIM).

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER LOOKLEN
      CALL ITOA(LOOKLEN(INPUT_PATH, 40), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER LOOKLEN
      CALL ITOA(LOOKLEN(INPUT_PATH, 50), ANSWER, ANSWER_LEN)
      END

C     LOOKLEN — run the look-and-say transform ITERS times, returning the
C     final length. A and B are static ping-pong buffers.
      INTEGER FUNCTION LOOKLEN(INPUT_PATH, ITERS)
      CHARACTER*(*) INPUT_PATH
      INTEGER ITERS
      INTEGER MAXL
      PARAMETER (MAXL = 8000000)
      CHARACTER A(MAXL), B(MAXL)
      INTEGER LA, LB, IT, I, J, RUN, IOS, K, D
      CHARACTER CUR
      CHARACTER*256 LINE
      CHARACTER*16 CNTBUF
      SAVE A, B

      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
      READ(10, '(A)', IOSTAT=IOS) LINE
      CLOSE(10)
      LA = LEN_TRIM(LINE)
      DO 10 I = 1, LA
        A(I) = LINE(I:I)
10    CONTINUE

      DO 50 IT = 1, ITERS
        LB = 0
        I = 1
20      IF (I .LE. LA) THEN
          CUR = A(I)
          RUN = 1
30        IF (I + RUN .LE. LA) THEN
            IF (A(I+RUN) .EQ. CUR) THEN
              RUN = RUN + 1
              GOTO 30
            END IF
          END IF
C         Emit the run count (as decimal digits) then the digit itself.
          WRITE(CNTBUF, '(I0)') RUN
          K = LEN_TRIM(CNTBUF)
          DO 40 J = 1, K
            LB = LB + 1
            B(LB) = CNTBUF(J:J)
40        CONTINUE
          LB = LB + 1
          B(LB) = CUR
          I = I + RUN
          GOTO 20
        END IF
C       Swap: B becomes the new A. Copy back (buffers are distinct arrays).
        DO 45 I = 1, LB
          A(I) = B(I)
45      CONTINUE
        LA = LB
50    CONTINUE

      LOOKLEN = LA
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
