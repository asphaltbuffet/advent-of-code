C Solution for Advent of Code 2015 day 16.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Each Sue lists a subset of compound readings; find the one consistent
C with the fixed MFCSAM target. Part 1 matches exactly. Part 2 treats
C cats/trees as "greater than" and pomeranians/goldfish as "fewer than".

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER FINDSUE
      CALL ITOA(FINDSUE(INPUT_PATH, 1), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER FINDSUE
      CALL ITOA(FINDSUE(INPUT_PATH, 2), ANSWER, ANSWER_LEN)
      END

C     FINDSUE — number of the first Sue whose every reading is consistent
C     with the target. PART selects exact (1) or ranged (2) comparisons.
      INTEGER FUNCTION FINDSUE(INPUT_PATH, PART)
      CHARACTER*(*) INPUT_PATH
      INTEGER PART
      CHARACTER*256 LINE
      CHARACTER*32 W(40)
      INTEGER IOS, NW, NUM, I, J, WANT, MODE, TVAL
      CHARACTER*32 COMP
      LOGICAL MATCH, OK

      NUM = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 100
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 20
      NUM = NUM + 1
      CALL SPLIT16(LINE, W, NW)
C     W(1)='Sue' W(2)='<num>:'. Compound/value pairs start at W(3):
C     "<comp>: <value>," — comp word ends in ':', value may end in ','.
      MATCH = .TRUE.
      I = 3
30    IF (I .LT. NW) THEN
        COMP = W(I)
        CALL STRIP(COMP, ':')
        READ(W(I+1), *) WANT
        CALL TARGET16(COMP, TVAL, MODE)
        IF (PART .EQ. 1) MODE = 0
        IF (MODE .EQ. 1) THEN
          OK = WANT .GT. TVAL
        ELSE IF (MODE .EQ. 2) THEN
          OK = WANT .LT. TVAL
        ELSE
          OK = WANT .EQ. TVAL
        END IF
        IF (.NOT. OK) MATCH = .FALSE.
        I = I + 2
        GOTO 30
      END IF
      IF (MATCH) THEN
        FINDSUE = NUM
        CLOSE(10)
        RETURN
      END IF
      GOTO 20
100   CLOSE(10)
      FINDSUE = -1
      END

C     TARGET16 — the MFCSAM reading for a compound: its value TVAL and a
C     comparison MODE (0 exact, 1 greater-than, 2 fewer-than) for part 2.
      SUBROUTINE TARGET16(COMP, TVAL, MODE)
      CHARACTER*32 COMP
      INTEGER TVAL, MODE

      MODE = 0
      IF (COMP .EQ. 'children') THEN
        TVAL = 3
      ELSE IF (COMP .EQ. 'cats') THEN
        TVAL = 7
        MODE = 1
      ELSE IF (COMP .EQ. 'samoyeds') THEN
        TVAL = 2
      ELSE IF (COMP .EQ. 'pomeranians') THEN
        TVAL = 3
        MODE = 2
      ELSE IF (COMP .EQ. 'akitas') THEN
        TVAL = 0
      ELSE IF (COMP .EQ. 'vizslas') THEN
        TVAL = 0
      ELSE IF (COMP .EQ. 'goldfish') THEN
        TVAL = 5
        MODE = 2
      ELSE IF (COMP .EQ. 'trees') THEN
        TVAL = 3
        MODE = 1
      ELSE IF (COMP .EQ. 'cars') THEN
        TVAL = 2
      ELSE IF (COMP .EQ. 'perfumes') THEN
        TVAL = 1
      ELSE
        TVAL = -1
      END IF
      END

C     STRIP — remove a trailing character CH from S (if present).
      SUBROUTINE STRIP(S, CH)
      CHARACTER*32 S
      CHARACTER CH
      INTEGER L
      L = LEN_TRIM(S)
      IF (L .GT. 0 .AND. S(L:L) .EQ. CH) S(L:L) = ' '
      END

C     SPLIT16 — break LINE into space-delimited words, stripping a trailing
C     comma from each (commas separate the compound:value pairs).
      SUBROUTINE SPLIT16(LINE, W, NW)
      CHARACTER*256 LINE
      CHARACTER*32 W(40)
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
      CALL STRIP(W(NW), ',')
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
