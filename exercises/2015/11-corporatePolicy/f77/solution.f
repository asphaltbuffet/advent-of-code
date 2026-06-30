C Solution for Advent of Code 2015 day 11.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Find the next valid password after the input: an increasing straight of
C three letters, no i/o/l, and two distinct non-overlapping pairs.
C Passwords are held as digit arrays (0..25) and incremented base-26.
C Part 2 is the next valid password after part 1's.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER P(8)

      CALL LOAD(INPUT_PATH, P)
      CALL NEXTPW(P)
      CALL PWOUT(P, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER P(8)

      CALL LOAD(INPUT_PATH, P)
      CALL NEXTPW(P)
      CALL INCR(P)
      CALL NEXTPW(P)
      CALL PWOUT(P, ANSWER, ANSWER_LEN)
      END

C     NEXTPW — advance P to the next password satisfying all rules.
      SUBROUTINE NEXTPW(P)
      INTEGER P(8)
      LOGICAL VALID

10    CALL INCR(P)
      IF (.NOT. VALID(P)) GOTO 10
      END

C     INCR — base-26 increment with carry. After incrementing, if any letter
C     is forbidden (i=8, o=14, l=11), bump it and zero everything to its
C     right, skipping the whole doomed range.
      SUBROUTINE INCR(P)
      INTEGER P(8)
      INTEGER I

      I = 8
10    P(I) = P(I) + 1
      IF (P(I) .GT. 25) THEN
        P(I) = 0
        I = I - 1
        IF (I .GE. 1) GOTO 10
      END IF

C     Forbidden-letter skip: leftmost i/o/l gets bumped, rest zeroed.
      DO 20 I = 1, 8
        IF (P(I) .EQ. 8 .OR. P(I) .EQ. 14 .OR. P(I) .EQ. 11) THEN
          P(I) = P(I) + 1
          IF (I .LT. 8) THEN
            DO 15 J = I + 1, 8
              P(J) = 0
15          CONTINUE
          END IF
          GOTO 30
        END IF
20    CONTINUE
30    CONTINUE
      END

C     VALID — true if P has a 3-letter straight and two distinct pairs.
C     (Forbidden letters are already excluded by INCR.)
      LOGICAL FUNCTION VALID(P)
      INTEGER P(8)
      INTEGER I, NPAIR, LAST
      LOGICAL STRAIGHT

      STRAIGHT = .FALSE.
      DO 10 I = 1, 6
        IF (P(I+1) .EQ. P(I) + 1 .AND. P(I+2) .EQ. P(I) + 2)
     &      STRAIGHT = .TRUE.
10    CONTINUE

C     Count non-overlapping pairs of distinct letters.
      NPAIR = 0
      LAST = -1
      I = 1
20    IF (I .LE. 7) THEN
        IF (P(I) .EQ. P(I+1) .AND. P(I) .NE. LAST) THEN
          NPAIR = NPAIR + 1
          LAST = P(I)
          I = I + 2
        ELSE
          I = I + 1
        END IF
        GOTO 20
      END IF

      VALID = STRAIGHT .AND. (NPAIR .GE. 2)
      END

C     LOAD — read the 8-letter password into digit values 0..25.
      SUBROUTINE LOAD(INPUT_PATH, P)
      CHARACTER*(*) INPUT_PATH
      INTEGER P(8)
      CHARACTER*256 LINE
      INTEGER IOS, I

      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
      READ(10, '(A)', IOSTAT=IOS) LINE
      CLOSE(10)
      DO 10 I = 1, 8
        P(I) = ICHAR(LINE(I:I)) - ICHAR('a')
10    CONTINUE
      END

C     PWOUT — render the digit array back to letters in the answer buffer.
      SUBROUTINE PWOUT(P, S, L)
      INTEGER P(8), L, I
      CHARACTER*256 S
      DO 10 I = 1, 8
        S(I:I) = CHAR(P(I) + ICHAR('a'))
10    CONTINUE
      L = 8
      END
