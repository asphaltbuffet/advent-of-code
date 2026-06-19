C Solution for Advent of Code 2015 day 13.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Seat guests around a circular table to maximise total happiness (each
C neighbour pair contributes both directions). Part 2 adds an indifferent
C "me". Names are interned to ids; seatings are brute-forced by permuting
C all but the first guest (factoring out rotations).

      BLOCK DATA DAY13INIT
      INTEGER NG, H, PERM, BEST
      CHARACTER*32 NAMES
      COMMON /D13/ H(20,20), PERM(20), NG, BEST
      COMMON /D13C/ NAMES(20)
      DATA NG /0/
      END

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER NG, H, PERM, BEST
      COMMON /D13/ H(20,20), PERM(20), NG, BEST

      CALL LOAD(INPUT_PATH)
      BEST = -2000000000
      CALL PERMUTE(2)
      CALL ITOA(BEST, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER NG, H, PERM, BEST
      COMMON /D13/ H(20,20), PERM(20), NG, BEST
      INTEGER I

      CALL LOAD(INPUT_PATH)
C     Add "me": indifferent to everyone (zero happiness both ways).
      DO 10 I = 1, NG
        H(NG+1, I) = 0
        H(I, NG+1) = 0
10    CONTINUE
      NG = NG + 1
      PERM(NG) = NG

      BEST = -2000000000
      CALL PERMUTE(2)
      CALL ITOA(BEST, ANSWER, ANSWER_LEN)
      END

C     PERMUTE — permute PERM(K..NG) in place; the first guest is fixed
C     (recursion starts at 2) to factor out circular rotations. Each full
C     seating sums happiness for both directions around the cycle.
      RECURSIVE SUBROUTINE PERMUTE(K)
      INTEGER K
      INTEGER NG, H, PERM, BEST
      COMMON /D13/ H(20,20), PERM(20), NG, BEST
      INTEGER I, T, TOTAL, A, B

      IF (K .EQ. NG + 1) THEN
        TOTAL = 0
        DO 10 I = 1, NG
          A = PERM(I)
          IF (I .EQ. NG) THEN
            B = PERM(1)
          ELSE
            B = PERM(I+1)
          END IF
          TOTAL = TOTAL + H(A,B) + H(B,A)
10      CONTINUE
        IF (TOTAL .GT. BEST) BEST = TOTAL
        RETURN
      END IF
      DO 20 I = K, NG
        T = PERM(K)
        PERM(K) = PERM(I)
        PERM(I) = T
        CALL PERMUTE(K + 1)
        T = PERM(K)
        PERM(K) = PERM(I)
        PERM(I) = T
20    CONTINUE
      END

      INTEGER FUNCTION INTERN(NAME)
      CHARACTER*32 NAME
      INTEGER NG, H, PERM, BEST
      CHARACTER*32 NAMES
      COMMON /D13/ H(20,20), PERM(20), NG, BEST
      COMMON /D13C/ NAMES(20)
      INTEGER I

      DO 10 I = 1, NG
        IF (NAMES(I) .EQ. NAME) THEN
          INTERN = I
          RETURN
        END IF
10    CONTINUE
      NG = NG + 1
      NAMES(NG) = NAME
      PERM(NG) = NG
      INTERN = NG
      END

C     LOAD — parse "<A> would gain|lose <N> happiness units by sitting next
C     to <B>." Happiness is directional: H(A,B) = +/- N.
      SUBROUTINE LOAD(INPUT_PATH)
      CHARACTER*(*) INPUT_PATH
      INTEGER NG, H, PERM, BEST
      COMMON /D13/ H(20,20), PERM(20), NG, BEST
      CHARACTER*256 LINE
      CHARACTER*32 PA, PB, SIGN
      INTEGER IOS, N, A, B, INTERN

      NG = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 100
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 20
      CALL PARSE13(LINE, PA, SIGN, N, PB)
      IF (SIGN .EQ. 'lose') N = -N
      A = INTERN(PA)
      B = INTERN(PB)
      H(A,B) = N
      GOTO 20
100   CLOSE(10)
      END

C     PARSE13 — split the line into words and pull the two people, the
C     gain/lose sign, and the amount. The last word carries a trailing '.'.
      SUBROUTINE PARSE13(LINE, PA, SIGN, N, PB)
      CHARACTER*256 LINE
      CHARACTER*32 PA, SIGN, PB
      INTEGER N
      CHARACTER*32 W(16)
      INTEGER NW, I, P, L

      CALL SPLIT(LINE, W, NW)
      PA = W(1)
      SIGN = W(3)
      READ(W(4), *) N
      PB = W(NW)
C     Strip a trailing period from the last name.
      L = LEN_TRIM(PB)
      IF (L .GT. 0 .AND. PB(L:L) .EQ. '.') PB(L:L) = ' '
      END

C     SPLIT — break LINE into space-delimited words W(1..NW).
      SUBROUTINE SPLIT(LINE, W, NW)
      CHARACTER*256 LINE
      CHARACTER*32 W(16)
      INTEGER NW, P, L, I

      NW = 0
      P = 1
      L = LEN_TRIM(LINE)
10    IF (P .LE. L .AND. LINE(P:P) .EQ. ' ') THEN
        P = P + 1
        GOTO 10
      END IF
      IF (P .GT. L) RETURN
      NW = NW + 1
      W(NW) = ' '
      I = 1
20    IF (P .LE. L .AND. LINE(P:P) .NE. ' ') THEN
        IF (I .LE. 32) W(NW)(I:I) = LINE(P:P)
        I = I + 1
        P = P + 1
        GOTO 20
      END IF
      GOTO 10
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
