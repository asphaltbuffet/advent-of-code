C Solution for Advent of Code 2015 day 7.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C A circuit of 16-bit wires fed by AND/OR/NOT/LSHIFT/RSHIFT gates or direct
C signals. Wire names (1-2 lowercase letters) are perfect-hashed to an int
C id, so a dense array replaces a hash map. Wire 'a' is evaluated by
C memoised recursion. Part 2 feeds part 1's 'a' back into 'b' and recomputes.

C     Per-wire tables, indexed by perfect-hash id (1..729):
C       OPC  operation: 0 set, 1 AND, 2 OR, 3 LSHIFT, 4 RSHIFT, 5 NOT
C       L1/IL1  first operand value and whether it is a literal
C       L2/IL2  second operand (or shift amount)
C       DONE/VAL  memoisation flag and cached signal
      BLOCK DATA DAY7INIT
      INTEGER OPC, L1, L2, IL1, IL2, DONE, VAL
      COMMON /D7/ OPC(729), L1(729), L2(729), IL1(729), IL2(729),
     &            DONE(729), VAL(729)
      END

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER WID, EVAL

      CALL LOAD(INPUT_PATH)
      CALL ITOA(EVAL(WID('a')), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER OPC, L1, L2, IL1, IL2, DONE, VAL
      COMMON /D7/ OPC(729), L1(729), L2(729), IL1(729), IL2(729),
     &            DONE(729), VAL(729)
      INTEGER WID, EVAL, AVAL, B, I

      CALL LOAD(INPUT_PATH)
      AVAL = EVAL(WID('a'))
C     Override wire b with a's signal, clear the memo, recompute a.
      B = WID('b')
      OPC(B) = 0
      IL1(B) = 1
      L1(B) = AVAL
      DO 10 I = 1, 729
        DONE(I) = 0
10    CONTINUE
      CALL ITOA(EVAL(WID('a')), ANSWER, ANSWER_LEN)
      END

C     WID — perfect hash of a 1-2 char lowercase wire name to 1..729.
      INTEGER FUNCTION WID(NAME)
      CHARACTER*(*) NAME
      INTEGER C1, C2

      C1 = ICHAR(NAME(1:1)) - ICHAR('a') + 1
      IF (LEN(NAME) .GE. 2 .AND. NAME(2:2) .NE. ' ') THEN
        C2 = ICHAR(NAME(2:2)) - ICHAR('a') + 1
      ELSE
        C2 = 0
      END IF
      WID = C1 * 27 + C2 + 1
      END

C     EVAL — memoised recursive evaluation of wire W's 16-bit signal.
      RECURSIVE INTEGER FUNCTION EVAL(W)
      INTEGER W
      INTEGER OPC, L1, L2, IL1, IL2, DONE, VAL
      COMMON /D7/ OPC(729), L1(729), L2(729), IL1(729), IL2(729),
     &            DONE(729), VAL(729)
      INTEGER A, B, R, OPRND

      IF (DONE(W) .EQ. 1) THEN
        EVAL = VAL(W)
        RETURN
      END IF

      IF (OPC(W) .EQ. 0) THEN
        R = OPRND(IL1(W), L1(W))
      ELSE IF (OPC(W) .EQ. 5) THEN
        R = 65535 - OPRND(IL1(W), L1(W))
      ELSE
        A = OPRND(IL1(W), L1(W))
        B = OPRND(IL2(W), L2(W))
        IF (OPC(W) .EQ. 1) THEN
          R = IAND(A, B)
        ELSE IF (OPC(W) .EQ. 2) THEN
          R = IOR(A, B)
        ELSE IF (OPC(W) .EQ. 3) THEN
          R = IAND(ISHFT(A, B), 65535)
        ELSE
          R = ISHFT(A, -B)
        END IF
      END IF

      R = IAND(R, 65535)
      VAL(W) = R
      DONE(W) = 1
      EVAL = R
      END

C     OPRND — resolve an operand: a literal returns its value, otherwise it
C     is a wire id to evaluate.
      INTEGER FUNCTION OPRND(ISLIT, V)
      INTEGER ISLIT, V, EVAL
      IF (ISLIT .EQ. 1) THEN
        OPRND = V
      ELSE
        OPRND = EVAL(V)
      END IF
      END

C     LOAD — parse the circuit. Each line is "<expr> -> <wire>".
      SUBROUTINE LOAD(INPUT_PATH)
      CHARACTER*(*) INPUT_PATH
      INTEGER OPC, L1, L2, IL1, IL2, DONE, VAL
      COMMON /D7/ OPC(729), L1(729), L2(729), IL1(729), IL2(729),
     &            DONE(729), VAL(729)
      CHARACTER*256 LINE
      CHARACTER*32 W(8)
      INTEGER IOS, NW, T, WID, I

      DO 5 I = 1, 729
        DONE(I) = 0
        OPC(I) = 0
        IL1(I) = 1
        L1(I) = 0
5     CONTINUE

      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 100
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 20
      CALL SPLIT7(LINE, W, NW)
C     Forms (the "->" token is dropped by SPLIT7, target is the last word):
C       x -> t            (NW=2)  set
C       NOT x -> t        (NW=3)  not
C       x OP y -> t       (NW=4)  binary
      T = WID(W(NW))
      IF (NW .EQ. 2) THEN
        OPC(T) = 0
        CALL OPND(W(1), IL1(T), L1(T))
      ELSE IF (NW .EQ. 3) THEN
        OPC(T) = 5
        CALL OPND(W(2), IL1(T), L1(T))
      ELSE
        CALL OPND(W(1), IL1(T), L1(T))
        CALL OPND(W(3), IL2(T), L2(T))
        IF (W(2) .EQ. 'AND') THEN
          OPC(T) = 1
        ELSE IF (W(2) .EQ. 'OR') THEN
          OPC(T) = 2
        ELSE IF (W(2) .EQ. 'LSHIFT') THEN
          OPC(T) = 3
        ELSE
          OPC(T) = 4
        END IF
      END IF
      GOTO 20
100   CLOSE(10)
      END

C     OPND — classify a token as a numeric literal or a wire reference,
C     setting ISLIT (1 literal, 0 wire) and V (the value or wire id).
      SUBROUTINE OPND(TOK, ISLIT, V)
      CHARACTER*32 TOK
      INTEGER ISLIT, V, WID
      CHARACTER C

      C = TOK(1:1)
      IF (C .GE. '0' .AND. C .LE. '9') THEN
        ISLIT = 1
        READ(TOK, *) V
      ELSE
        ISLIT = 0
        V = WID(TOK)
      END IF
      END

C     SPLIT7 — break LINE into space-delimited words, dropping the "->".
      SUBROUTINE SPLIT7(LINE, W, NW)
      CHARACTER*256 LINE
      CHARACTER*32 W(8)
      INTEGER NW, P, L, I
      CHARACTER*32 TOK

      NW = 0
      P = 1
      L = LEN_TRIM(LINE)
10    IF (P .LE. L .AND. LINE(P:P) .EQ. ' ') THEN
        P = P + 1
        GOTO 10
      END IF
      IF (P .GT. L) RETURN
      TOK = ' '
      I = 1
20    IF (P .LE. L .AND. LINE(P:P) .NE. ' ') THEN
        IF (I .LE. 32) TOK(I:I) = LINE(P:P)
        I = I + 1
        P = P + 1
        GOTO 20
      END IF
      IF (TOK .NE. '->') THEN
        NW = NW + 1
        W(NW) = TOK
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
