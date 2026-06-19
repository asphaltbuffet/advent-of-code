C Solution for Advent of Code 2015 day 2.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Each line is LxWxH. Paper = surface area + smallest face. Ribbon =
C smallest face perimeter + volume.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER L, W, H, IOS, TOTAL, A, B, C, SMALL
      CHARACTER*256 LINE

      TOTAL = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
10    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 20
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 10
      CALL DIGITS(LINE)
      READ(LINE, *) L, W, H
      A = L * W
      B = W * H
      C = H * L
      SMALL = A
      IF (B .LT. SMALL) SMALL = B
      IF (C .LT. SMALL) SMALL = C
      TOTAL = TOTAL + 2*A + 2*B + 2*C + SMALL
      GOTO 10
20    CLOSE(10)
      CALL ITOA(TOTAL, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER L, W, H, IOS, TOTAL, P1, P2, P3, SMALL
      CHARACTER*256 LINE

      TOTAL = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
10    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 20
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 10
      CALL DIGITS(LINE)
      READ(LINE, *) L, W, H
      P1 = 2 * (L + W)
      P2 = 2 * (W + H)
      P3 = 2 * (H + L)
      SMALL = P1
      IF (P2 .LT. SMALL) SMALL = P2
      IF (P3 .LT. SMALL) SMALL = P3
      TOTAL = TOTAL + SMALL + L * W * H
      GOTO 10
20    CLOSE(10)
      CALL ITOA(TOTAL, ANSWER, ANSWER_LEN)
      END

C     DIGITS — blank every non-digit character so list-directed READ can pull
C     the integers regardless of the 'x' separators.
      SUBROUTINE DIGITS(LINE)
      CHARACTER*256 LINE
      INTEGER I
      DO 10 I = 1, 256
        IF (LINE(I:I) .LT. '0' .OR. LINE(I:I) .GT. '9')
     &      LINE(I:I) = ' '
10    CONTINUE
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
