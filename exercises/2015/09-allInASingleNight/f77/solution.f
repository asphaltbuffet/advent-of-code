C Solution for Advent of Code 2015 day 9.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Shortest (part 1) and longest (part 2) route visiting every city once.
C City names are interned to integer ids (linear scan stands in for a hash
C map); routes are brute-forced by recursive permutation (N! <= 8!).

C     Shared state: NAMES/NC the interned cities; DM the distance matrix;
C     PERM the working permutation; BESTMIN/BESTMAX the route extremes.
      BLOCK DATA DAY9INIT
      INTEGER NC, DM, PERM, BESTMIN, BESTMAX
      CHARACTER*32 NAMES
      COMMON /D9/ DM(20,20), PERM(20), NC, BESTMIN, BESTMAX
      COMMON /D9C/ NAMES(20)
      DATA NC /0/
      END

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER NC, DM, PERM, BESTMIN, BESTMAX
      COMMON /D9/ DM(20,20), PERM(20), NC, BESTMIN, BESTMAX

      CALL LOAD(INPUT_PATH)
      BESTMIN = 2000000000
      BESTMAX = -1
      CALL PERMUTE(1)
      CALL ITOA(BESTMIN, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER NC, DM, PERM, BESTMIN, BESTMAX
      COMMON /D9/ DM(20,20), PERM(20), NC, BESTMIN, BESTMAX

      CALL LOAD(INPUT_PATH)
      BESTMIN = 2000000000
      BESTMAX = -1
      CALL PERMUTE(1)
      CALL ITOA(BESTMAX, ANSWER, ANSWER_LEN)
      END

C     PERMUTE — recursively generate all orderings of PERM(K..NC) in place,
C     scoring each complete route against BESTMIN/BESTMAX.
      RECURSIVE SUBROUTINE PERMUTE(K)
      INTEGER K
      INTEGER NC, DM, PERM, BESTMIN, BESTMAX
      COMMON /D9/ DM(20,20), PERM(20), NC, BESTMIN, BESTMAX
      INTEGER I, T, TOTAL

      IF (K .EQ. NC) THEN
        TOTAL = 0
        DO 10 I = 1, NC - 1
          TOTAL = TOTAL + DM(PERM(I), PERM(I+1))
10      CONTINUE
        IF (TOTAL .LT. BESTMIN) BESTMIN = TOTAL
        IF (TOTAL .GT. BESTMAX) BESTMAX = TOTAL
        RETURN
      END IF
      DO 20 I = K, NC
        T = PERM(K)
        PERM(K) = PERM(I)
        PERM(I) = T
        CALL PERMUTE(K + 1)
        T = PERM(K)
        PERM(K) = PERM(I)
        PERM(I) = T
20    CONTINUE
      END

C     INTERN — return the id of city NAME, adding it if new.
      INTEGER FUNCTION INTERN(NAME)
      CHARACTER*32 NAME
      INTEGER NC, DM, PERM, BESTMIN, BESTMAX
      CHARACTER*32 NAMES
      COMMON /D9/ DM(20,20), PERM(20), NC, BESTMIN, BESTMAX
      COMMON /D9C/ NAMES(20)
      INTEGER I

      DO 10 I = 1, NC
        IF (NAMES(I) .EQ. NAME) THEN
          INTERN = I
          RETURN
        END IF
10    CONTINUE
      NC = NC + 1
      NAMES(NC) = NAME
      PERM(NC) = NC
      INTERN = NC
      END

C     LOAD — parse "<src> to <dst> = <dist>" lines into the distance matrix.
      SUBROUTINE LOAD(INPUT_PATH)
      CHARACTER*(*) INPUT_PATH
      INTEGER NC, DM, PERM, BESTMIN, BESTMAX
      COMMON /D9/ DM(20,20), PERM(20), NC, BESTMIN, BESTMAX
      CHARACTER*256 LINE
      CHARACTER*32 W1, W2
      INTEGER IOS, DIST, A, B, INTERN

      NC = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 100
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 20
      CALL PARSE9(LINE, W1, W2, DIST)
      A = INTERN(W1)
      B = INTERN(W2)
      DM(A,B) = DIST
      DM(B,A) = DIST
      GOTO 20
100   CLOSE(10)
      END

C     PARSE9 — pull source name, destination name, and distance from a line
C     of the form "<src> to <dst> = <dist>".
      SUBROUTINE PARSE9(LINE, SRC, DST, DIST)
      CHARACTER*256 LINE
      CHARACTER*32 SRC, DST
      INTEGER DIST
      CHARACTER*32 W(8)
      INTEGER NW, I, P, L

      NW = 0
      P = 1
      L = LEN_TRIM(LINE)
C     Split on spaces into words W(1..NW).
10    CONTINUE
      DO 20 WHILE (P .LE. L .AND. LINE(P:P) .EQ. ' ')
        P = P + 1
20    CONTINUE
      IF (P .GT. L) GOTO 40
      NW = NW + 1
      W(NW) = ' '
      I = 1
30    IF (P .LE. L .AND. LINE(P:P) .NE. ' ') THEN
        W(NW)(I:I) = LINE(P:P)
        I = I + 1
        P = P + 1
        GOTO 30
      END IF
      GOTO 10
40    CONTINUE
C     "<src>=W1 to=W2 <dst>=W3 ==W4 <dist>=W5"
      SRC = W(1)
      DST = W(3)
      READ(W(5), *) DIST
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
