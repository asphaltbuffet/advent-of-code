C Solution for Advent of Code 2015 day 21.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Buy gear (1 weapon, 0-1 armor, 0-2 distinct rings) from the fixed shop and
C fight the boss. Part 1: least gold that still wins. Part 2: most gold that
C still loses. The player (100 HP) attacks first, so the fight has a closed
C form: compare ceil(targetHP / effectiveDamage) for each side.

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER SHOP
      CALL ITOA(SHOP(INPUT_PATH, 1), ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER SHOP
      CALL ITOA(SHOP(INPUT_PATH, 2), ANSWER, ANSWER_LEN)
      END

C     SHOP — enumerate every legal loadout; WANT 1 returns the cheapest win,
C     WANT 2 the priciest loss.
      INTEGER FUNCTION SHOP(INPUT_PATH, WANT)
      CHARACTER*(*) INPUT_PATH
      INTEGER WANT
      INTEGER WC(5), WD(5), AC(6), AAR(6), RC(6), RD(6), RAR(6)
      INTEGER BHP, BDMG, BARM, W, A, I, J, COST, DMG, ARM
      INTEGER BEST, WORST
      LOGICAL WINS
C     Weapons: cost, damage.
      DATA WC /8, 10, 25, 40, 74/
      DATA WD /4, 5, 6, 7, 8/
C     Armor (index 1 = none): cost, armor.
      DATA AC /0, 13, 31, 53, 75, 102/
      DATA AAR /0, 1, 2, 3, 4, 5/
C     Rings: cost, damage, armor.
      DATA RC /25, 50, 100, 20, 40, 80/
      DATA RD /1, 2, 3, 0, 0, 0/
      DATA RAR /0, 0, 0, 1, 2, 3/

      CALL READBOSS(INPUT_PATH, BHP, BDMG, BARM)

      BEST = -1
      WORST = 0
      DO 50 W = 1, 5
        DO 40 A = 1, 6
C         Ring sets: I=0 none; I>0 picks ring I, J>I adds a distinct second.
          DO 30 I = 0, 6
            DO 20 J = I, 6
              IF (I .EQ. 0 .AND. J .GT. 0) GOTO 20
              COST = WC(W) + AC(A)
              DMG = WD(W)
              ARM = AAR(A)
              IF (I .GT. 0) THEN
                COST = COST + RC(I)
                DMG = DMG + RD(I)
                ARM = ARM + RAR(I)
              END IF
              IF (J .GT. I) THEN
                COST = COST + RC(J)
                DMG = DMG + RD(J)
                ARM = ARM + RAR(J)
              END IF
              IF (WINS(DMG, ARM, BHP, BDMG, BARM)) THEN
                IF (BEST .LT. 0 .OR. COST .LT. BEST) BEST = COST
              ELSE
                IF (COST .GT. WORST) WORST = COST
              END IF
20          CONTINUE
30        CONTINUE
40      CONTINUE
50    CONTINUE

      IF (WANT .EQ. 1) THEN
        SHOP = BEST
      ELSE
        SHOP = WORST
      END IF
      END

C     WINS — true if the player (100 HP, first strike) beats the boss.
      LOGICAL FUNCTION WINS(PD, PA, BHP, BDMG, BARM)
      INTEGER PD, PA, BHP, BDMG, BARM
      INTEGER DB, DP, BT, PT

      DB = PD - BARM
      IF (DB .LT. 1) DB = 1
      DP = BDMG - PA
      IF (DP .LT. 1) DP = 1
      BT = (BHP + DB - 1) / DB
      PT = (100 + DP - 1) / DP
      WINS = BT .LE. PT
      END

C     READBOSS — parse the boss HP, damage, armor (one integer per line).
      SUBROUTINE READBOSS(INPUT_PATH, BHP, BDMG, BARM)
      CHARACTER*(*) INPUT_PATH
      INTEGER BHP, BDMG, BARM
      CHARACTER*64 LINE
      INTEGER IOS, K

C     Lines are "Hit Points: N", "Damage: N", "Armor: N" in that order; strip
C     non-digits then read the value.
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
      READ(10, '(A)') LINE
      CALL DIGITS(LINE)
      READ(LINE, *) BHP
      READ(10, '(A)') LINE
      CALL DIGITS(LINE)
      READ(LINE, *) BDMG
      READ(10, '(A)') LINE
      CALL DIGITS(LINE)
      READ(LINE, *) BARM
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
