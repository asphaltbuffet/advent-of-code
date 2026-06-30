C Solution for Advent of Code 2015 day 19.
C
C Implement PART_ONE and PART_TWO. Each receives the path to a file
C containing the puzzle input and must write the answer into ANSWER,
C setting ANSWER_LEN to the number of characters written.
C
C Fixed-form Fortran 77: columns 1-6 are label/continuation, 7-72 code.
C
C Part 1: count distinct molecules from one replacement (dedup via a
C linear-scanned string array, since F77 has no hash set). Part 2: greedily
C reduce the molecule back to "e", reshuffling rules (tiny LCG) on a dead
C end. String-search work uses INDEX/substr.

C     Shared rule table and molecule.
      BLOCK DATA DAY19INIT
      INTEGER NR
      CHARACTER*64 FR, TO
      CHARACTER*1024 MOL
      COMMON /D19/ NR
      COMMON /D19C/ FR(100), TO(100), MOL
      DATA NR /0/
      END

      SUBROUTINE PART_ONE(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER NR
      CHARACTER*64 FR, TO
      CHARACTER*1024 MOL
      COMMON /D19/ NR
      COMMON /D19C/ FR(100), TO(100), MOL
C     Distinct-molecule store: up to NMAX candidates of up to 512 chars.
      INTEGER NMAX
      PARAMETER (NMAX = 2000)
      CHARACTER*512 UNIQ(NMAX)
      INTEGER NU, I, J, K, FL, TL, ML, POS, EXISTS
      CHARACTER*512 CAND
      SAVE UNIQ

      CALL LOAD(INPUT_PATH)
      ML = LEN_TRIM(MOL)
      NU = 0
      DO 40 I = 1, NR
        FL = LEN_TRIM(FR(I))
        TL = LEN_TRIM(TO(I))
        DO 30 POS = 1, ML - FL + 1
          IF (MOL(POS:POS+FL-1) .EQ. FR(I)(1:FL)) THEN
            CAND = MOL(1:POS-1) // TO(I)(1:TL) // MOL(POS+FL:ML)
            EXISTS = 0
            DO 20 K = 1, NU
              IF (UNIQ(K) .EQ. CAND) THEN
                EXISTS = 1
                GOTO 25
              END IF
20          CONTINUE
25          IF (EXISTS .EQ. 0) THEN
              NU = NU + 1
              UNIQ(NU) = CAND
            END IF
          END IF
30      CONTINUE
40    CONTINUE
      CALL ITOA(NU, ANSWER, ANSWER_LEN)
      END

      SUBROUTINE PART_TWO(INPUT_PATH, ANSWER, ANSWER_LEN)
      CHARACTER*(*) INPUT_PATH
      CHARACTER*256 ANSWER
      INTEGER ANSWER_LEN
      INTEGER NR
      CHARACTER*64 FR, TO
      CHARACTER*1024 MOL
      COMMON /D19/ NR
      COMMON /D19C/ FR(100), TO(100), MOL
      INTEGER ORD(100), I, K, J, TL, STEPS, STUCK, APPLIED, IDX, TMP, CL
      CHARACTER*1024 CUR
      INTEGER*8 STATE

      CALL LOAD(INPUT_PATH)
      DO 5 I = 1, NR
        ORD(I) = I
5     CONTINUE
      STATE = 1

10    CONTINUE
      CUR = MOL
      STEPS = 0
      STUCK = 0
20    CL = LEN_TRIM(CUR)
      IF (.NOT. (CL .EQ. 1 .AND. CUR(1:1) .EQ. 'e')) THEN
        APPLIED = 0
        DO 30 K = 1, NR
          I = ORD(K)
          TL = LEN_TRIM(TO(I))
          J = INDEX(CUR(1:CL), TO(I)(1:TL))
          IF (J .GT. 0) THEN
            CUR = CUR(1:J-1) // FR(I)(1:LEN_TRIM(FR(I)))
     &            // CUR(J+TL:CL)
            STEPS = STEPS + 1
            APPLIED = 1
            GOTO 35
          END IF
30      CONTINUE
35      IF (APPLIED .EQ. 0) THEN
          STUCK = 1
          GOTO 50
        END IF
        GOTO 20
      END IF
50    CONTINUE
      IF (STUCK .EQ. 0) THEN
        CALL ITOA(STEPS, ANSWER, ANSWER_LEN)
        RETURN
      END IF
C     Dead end: Fisher-Yates shuffle ORD via a tiny LCG, then retry.
      DO 60 K = NR, 2, -1
        STATE = STATE * 6364136223846793005_8 + 1
        IDX = INT(MOD(ISHFT(STATE, -33), INT(K, 8))) + 1
        TMP = ORD(K)
        ORD(K) = ORD(IDX)
        ORD(IDX) = TMP
60    CONTINUE
      GOTO 10
      END

C     LOAD — parse the rules ("from => to") and the molecule (after a blank
C     line) into the shared common block.
      SUBROUTINE LOAD(INPUT_PATH)
      CHARACTER*(*) INPUT_PATH
      INTEGER NR
      CHARACTER*64 FR, TO
      CHARACTER*1024 MOL
      COMMON /D19/ NR
      COMMON /D19C/ FR(100), TO(100), MOL
      CHARACTER*1024 LINE
      INTEGER IOS, P, BLANK

      NR = 0
      MOL = ' '
      BLANK = 0
      OPEN(10, FILE=INPUT_PATH, STATUS='OLD')
20    READ(10, '(A)', IOSTAT=IOS) LINE
      IF (IOS .NE. 0) GOTO 100
      IF (LEN_TRIM(LINE) .EQ. 0) THEN
        BLANK = 1
        GOTO 20
      END IF
      IF (BLANK .EQ. 0) THEN
C       "from => to"
        P = INDEX(LINE, ' => ')
        NR = NR + 1
        FR(NR) = LINE(1:P-1)
        TO(NR) = LINE(P+4:LEN_TRIM(LINE))
      ELSE
        MOL = LINE(1:LEN_TRIM(LINE))
      END IF
      GOTO 20
100   CLOSE(10)
      END

      SUBROUTINE ITOA(N, S, L)
      INTEGER N, L
      CHARACTER*256 S
      CHARACTER*32 BUF
      WRITE(BUF, '(I0)') N
      S = BUF
      L = LEN_TRIM(BUF)
      END
