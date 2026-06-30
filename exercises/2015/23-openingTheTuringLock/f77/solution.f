C Solution for Advent of Code 2015 day 23.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C A small assembly interpreter with registers a and b. Six opcodes:
C hlf/tpl/inc act on a register; jmp/jie/jio branch. jio jumps if the
C register equals one (not if it is odd). Part 1 starts a=0, part 2 a=1;
C both report register b at halt (PC stepping outside the program).

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER RUN23
      CALL ITOA(RUN23(INPUT_PATH, 0), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER RUN23
      CALL ITOA(RUN23(INPUT_PATH, 1), ANSWER, ANSWER_LEN)
      END

C     RUN23 — load the program, run it with register a = STARTA, return b.
      INTEGER FUNCTION RUN23(INPUT_PATH, STARTA)
      CHARACTER*(*) INPUT_PATH
      INTEGER STARTA
C     Decoded program: OPC opcode (1 hlf,2 tpl,3 inc,4 jmp,5 jie,6 jio),
C     REGI register index (1=a,2=b), OFF jump offset.
      INTEGER OPC(1000), REGI(1000), OFF(1000), N
      INTEGER REGS(2), PC, R
      CHARACTER*64 LINE
      INTEGER IOS

      N = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 30
      IF (LEN_TRIM(LINE) .EQ. 0) GOTO 20
      N = N + 1
      CALL DECODE(LINE, OPC(N), REGI(N), OFF(N))
      GOTO 20
30    CLOSE(10)

      REGS(1) = STARTA
      REGS(2) = 0
      PC = 1
40    IF (PC .GE. 1 .AND. PC .LE. N) THEN
        R = REGI(PC)
        IF (OPC(PC) .EQ. 1) THEN
          REGS(R) = REGS(R) / 2
          PC = PC + 1
        ELSE IF (OPC(PC) .EQ. 2) THEN
          REGS(R) = REGS(R) * 3
          PC = PC + 1
        ELSE IF (OPC(PC) .EQ. 3) THEN
          REGS(R) = REGS(R) + 1
          PC = PC + 1
        ELSE IF (OPC(PC) .EQ. 4) THEN
          PC = PC + OFF(PC)
        ELSE IF (OPC(PC) .EQ. 5) THEN
          IF (MOD(REGS(R), 2) .EQ. 0) THEN
            PC = PC + OFF(PC)
          ELSE
            PC = PC + 1
          END IF
        ELSE
          IF (REGS(R) .EQ. 1) THEN
            PC = PC + OFF(PC)
          ELSE
            PC = PC + 1
          END IF
        END IF
        GOTO 40
      END IF

      RUN23 = REGS(2)
      END

C     DECODE — parse one instruction line into opcode, register, and offset.
      SUBROUTINE DECODE(LINE, OP, REG, OFF)
      CHARACTER*64 LINE
      INTEGER OP, REG, OFF
      CHARACTER*3 MNE
      CHARACTER RC
      INTEGER I

      MNE = LINE(1:3)
      OP = 0
      REG = 1
      OFF = 0
      IF (MNE .EQ. 'hlf') OP = 1
      IF (MNE .EQ. 'tpl') OP = 2
      IF (MNE .EQ. 'inc') OP = 3
      IF (MNE .EQ. 'jmp') OP = 4
      IF (MNE .EQ. 'jie') OP = 5
      IF (MNE .EQ. 'jio') OP = 6

C     Register operand (if any) is the first a/b after the mnemonic.
      IF (OP .NE. 4) THEN
        RC = LINE(5:5)
        IF (RC .EQ. 'b') THEN
          REG = 2
        ELSE
          REG = 1
        END IF
      END IF

C     Offset (jmp/jie/jio): blank everything but digits and sign, then read.
      IF (OP .GE. 4) THEN
        DO 10 I = 1, 64
          RC = LINE(I:I)
          IF (.NOT. ((RC .GE. '0' .AND. RC .LE. '9')
     &        .OR. RC .EQ. '+' .OR. RC .EQ. '-'))
     &        LINE(I:I) = ' '
10      CONTINUE
        READ(LINE, *) OFF
      END IF
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
