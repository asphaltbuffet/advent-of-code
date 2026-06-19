C Solution for Advent of Code 2015 day 22.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Wizard fight: cast spells (some instant, some timed effects) to defeat the
C boss for the least mana. Part 2 is hard mode (the player loses 1 HP at the
C start of each of their turns). A recursive branch-and-bound DFS finds the
C minimum; BEST/BDMG/HARD are shared, the per-fight state is passed by value.

      BLOCK DATA DAY22INIT
      INTEGER BEST, BDMG, HARD
      COMMON /D22/ BEST, BDMG, HARD
      END

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER RUN22
      CALL ITOA(RUN22(INPUT_PATH, 0), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER RUN22
      CALL ITOA(RUN22(INPUT_PATH, 1), ANSWER, ANSWER_LEN)
      END

C     RUN22 — parse the boss, run the DFS, return the least mana to win.
      INTEGER FUNCTION RUN22(INPUT_PATH, HMODE)
      CHARACTER*(*) INPUT_PATH
      INTEGER HMODE
      INTEGER BEST, BDMG, HARD
      COMMON /D22/ BEST, BDMG, HARD
      INTEGER BHP

      CALL READBOSS(INPUT_PATH, BHP, BDMG)
      HARD = HMODE
      BEST = 1000000000
      CALL SEARCH(50, 500, BHP, 0, 0, 0, 0)
      RUN22 = BEST
      END

C     SEARCH — one player+boss cycle from the given state, branching over each
C     castable spell. Arguments are copied to locals first (F77 passes by
C     reference, so the recursion must not mutate the caller's values).
      RECURSIVE SUBROUTINE SEARCH(PHP0, MANA0, BHP0, SH0, PO0, RE0, SP0)
      INTEGER PHP0, MANA0, BHP0, SH0, PO0, RE0, SP0
      INTEGER BEST, BDMG, HARD
      COMMON /D22/ BEST, BDMG, HARD
      INTEGER SCOST(5), SDMG(5), SHEAL(5), SSH(5), SPO(5), SRE(5)
      INTEGER PHP, MANA, BHP, SH, PO, RE, SP
      INTEGER I, NP, NM, NB, NSH, NPO, NRE, NSP, HIT, ARMOR
C     cost / dmg / heal / shield / poison / recharge
      DATA SCOST /53, 73, 113, 173, 229/
      DATA SDMG /4, 2, 0, 0, 0/
      DATA SHEAL /0, 2, 0, 0, 0/
      DATA SSH /0, 0, 6, 0, 0/
      DATA SPO /0, 0, 0, 6, 0/
      DATA SRE /0, 0, 0, 0, 5/

      SP = SP0
      IF (SP .GE. BEST) RETURN

      PHP = PHP0
      MANA = MANA0
      BHP = BHP0
      SH = SH0
      PO = PO0
      RE = RE0

C     --- Player turn ---
      IF (HARD .EQ. 1) THEN
        PHP = PHP - 1
        IF (PHP .LE. 0) RETURN
      END IF
      IF (PO .GT. 0) THEN
        BHP = BHP - 3
        PO = PO - 1
      END IF
      IF (RE .GT. 0) THEN
        MANA = MANA + 101
        RE = RE - 1
      END IF
      IF (SH .GT. 0) SH = SH - 1
      IF (BHP .LE. 0) THEN
        IF (SP .LT. BEST) BEST = SP
        RETURN
      END IF

      DO 100 I = 1, 5
        IF (SCOST(I) .GT. MANA) GOTO 100
        IF (SSH(I) .GT. 0 .AND. SH .GT. 0) GOTO 100
        IF (SPO(I) .GT. 0 .AND. PO .GT. 0) GOTO 100
        IF (SRE(I) .GT. 0 .AND. RE .GT. 0) GOTO 100

        NP = PHP + SHEAL(I)
        NM = MANA - SCOST(I)
        NB = BHP - SDMG(I)
        IF (SSH(I) .GT. 0) THEN
          NSH = SSH(I)
        ELSE
          NSH = SH
        END IF
        IF (SPO(I) .GT. 0) THEN
          NPO = SPO(I)
        ELSE
          NPO = PO
        END IF
        IF (SRE(I) .GT. 0) THEN
          NRE = SRE(I)
        ELSE
          NRE = RE
        END IF
        NSP = SP + SCOST(I)

        IF (NB .LE. 0) THEN
          IF (NSP .LT. BEST) BEST = NSP
          GOTO 100
        END IF

C       --- Boss turn ---
        IF (NPO .GT. 0) THEN
          NB = NB - 3
          NPO = NPO - 1
        END IF
        IF (NRE .GT. 0) THEN
          NM = NM + 101
          NRE = NRE - 1
        END IF
        IF (NSH .GT. 0) NSH = NSH - 1
        IF (NB .LE. 0) THEN
          IF (NSP .LT. BEST) BEST = NSP
          GOTO 100
        END IF
        IF (NSH .GT. 0) THEN
          ARMOR = 7
        ELSE
          ARMOR = 0
        END IF
        HIT = BDMG - ARMOR
        IF (HIT .LT. 1) HIT = 1
        NP = NP - HIT
        IF (NP .LE. 0) GOTO 100

        CALL SEARCH(NP, NM, NB, NSH, NPO, NRE, NSP)
100   CONTINUE
      END

C     READBOSS — parse boss HP and damage (one integer per line).
      SUBROUTINE READBOSS(INPUT_PATH, BHP, BDMG)
      CHARACTER*(*) INPUT_PATH
      INTEGER BHP, BDMG
      CHARACTER*64 LINE

      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
      READ(10, '(A)') LINE
      CALL DIGITS(LINE)
      READ(LINE, *) BHP
      READ(10, '(A)') LINE
      CALL DIGITS(LINE)
      READ(LINE, *) BDMG
      CLOSE(10)
      END

      SUBROUTINE DIGITS(LINE)
      CHARACTER*64 LINE
      INTEGER I
      DO 10 I = 1, 64
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
