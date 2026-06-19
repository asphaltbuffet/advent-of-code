C Solution for Advent of Code 2015 day 24.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Split the packages into GROUPS equal-weight piles. The first pile must
C have the fewest packages, then the lowest quantum entanglement (product
C of its weights). Search by increasing first-pile size; the QE can reach
C ~1e10, so it is held in INTEGER*8.

      BLOCK DATA DAY24INIT
      INTEGER N
      INTEGER W
      INTEGER*8 BEST
      COMMON /D24/ W(64), N, BEST
      END

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER*8 MINQE
      CALL ITOA8(MINQE(INPUT_PATH, 3), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER*8 MINQE
      CALL ITOA8(MINQE(INPUT_PATH, 4), ANSWER, ANSWER_LEN)
      END

C     MINQE — load the weights, then for each increasing first-pile size find
C     the lowest quantum entanglement of a subset summing to total/GROUPS.
C     Returns as soon as a size yields any valid subset.
      INTEGER*8 FUNCTION MINQE(INPUT_PATH, GROUPS)
      CHARACTER*(*) INPUT_PATH
      INTEGER GROUPS
      INTEGER N
      INTEGER W
      INTEGER*8 BEST
      COMMON /D24/ W(64), N, BEST
      INTEGER TOTAL, TARGET, SIZE, IOS, V

      N = 0
      TOTAL = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, *, IOSTAT=IOS) V
      IF (IOS .NE. 0) GOTO 30
      N = N + 1
      W(N) = V
      TOTAL = TOTAL + V
      GOTO 20
30    CLOSE(10)

      TARGET = TOTAL / GROUPS

      DO 40 SIZE = 1, N
        BEST = -1
        CALL PICK(1, TARGET, SIZE, 1_8)
        IF (BEST .GE. 0) THEN
          MINQE = BEST
          RETURN
        END IF
40    CONTINUE
      MINQE = -1
      END

C     PICK — choose COUNT more weights from index START; track the minimum
C     product QE of selections that exactly reach REMAINING.
      RECURSIVE SUBROUTINE PICK(START, REMAINING, COUNT, QE)
      INTEGER START, REMAINING, COUNT
      INTEGER*8 QE
      INTEGER N
      INTEGER W
      INTEGER*8 BEST
      COMMON /D24/ W(64), N, BEST
      INTEGER I

      IF (COUNT .EQ. 0) THEN
        IF (REMAINING .EQ. 0 .AND. (BEST .LT. 0 .OR. QE .LT. BEST))
     &      BEST = QE
        RETURN
      END IF
      DO 10 I = START, N - COUNT + 1
        IF (W(I) .LE. REMAINING) THEN
          CALL PICK(I + 1, REMAINING - W(I), COUNT - 1,
     &              QE * INT(W(I), 8))
        END IF
10    CONTINUE
      END

      SUBROUTINE ITOA8(N, S, L)
      INTEGER*8 N
      INTEGER L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
