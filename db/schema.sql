-- CreateEnum
CREATE TYPE "auth_type_t" AS ENUM ('Google', 'Github', 'Email');

-- CreateTable
CREATE TABLE "User" (
    "id" SERIAL NOT NULL,
    "username" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "picture" TEXT NOT NULL,
    "auth_type" "auth_type_t" NOT NULL,

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Image" (
    "id" SERIAL NOT NULL,
    "associated_user" INTEGER NOT NULL,
    "filename" TEXT NOT NULL,
    "associated_post" INTEGER NOT NULL,
    "is_uploaded" BOOLEAN NOT NULL,

    CONSTRAINT "Image_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "LessonPost" (
    "id" SERIAL NOT NULL,
    "created_by" TEXT NOT NULL,
    "body" TEXT NOT NULL,
    "user_id" INTEGER NOT NULL,
    "status" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL,
    "is_public" BOOLEAN NOT NULL,
    "associated_course" INTEGER NOT NULL,

    CONSTRAINT "LessonPost_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "OverrideLessonPostVisibility" (
    "id" SERIAL NOT NULL,
    "allowed_user" INTEGER NOT NULL,
    "allowed_post_id" INTEGER NOT NULL,
    "allowed" BOOLEAN NOT NULL,

    CONSTRAINT "OverrideLessonPostVisibility_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "OverrideLessonPostPermission" (
    "id" SERIAL NOT NULL,
    "allowed_editor" INTEGER NOT NULL,
    "allowed_post_id" INTEGER NOT NULL,
    "allowed" BOOLEAN NOT NULL,

    CONSTRAINT "OverrideLessonPostPermission_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Course" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "language" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "is_public" BOOLEAN NOT NULL,

    CONSTRAINT "Course_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "CourseInstructorMapping" (
    "id" SERIAL NOT NULL,
    "instructor" INTEGER NOT NULL,
    "course_id" INTEGER NOT NULL,

    CONSTRAINT "CourseInstructorMapping_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "CourseLearnerMapping" (
    "id" SERIAL NOT NULL,
    "learner" INTEGER NOT NULL,
    "course_id" INTEGER NOT NULL,

    CONSTRAINT "CourseLearnerMapping_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "LiveClass" (
    "id" SERIAL NOT NULL,
    "start_time" TIMESTAMP(3) NOT NULL,
    "associated_course" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "reminder_message" TEXT NOT NULL,
    "is_public" BOOLEAN NOT NULL,

    CONSTRAINT "LiveClass_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "OverrideLiveClassVisibility" (
    "id" SERIAL NOT NULL,
    "learner" INTEGER NOT NULL,
    "live_class_id" INTEGER NOT NULL,

    CONSTRAINT "OverrideLiveClassVisibility_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "OverrideLiveClassAdminPermission" (
    "id" SERIAL NOT NULL,
    "instructor" INTEGER NOT NULL,
    "live_class_id" INTEGER NOT NULL,

    CONSTRAINT "OverrideLiveClassAdminPermission_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "User_username_key" ON "User"("username");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "Image_filename_key" ON "Image"("filename");

-- AddForeignKey
ALTER TABLE "Image" ADD CONSTRAINT "Image_associated_user_fkey" FOREIGN KEY ("associated_user") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Image" ADD CONSTRAINT "Image_associated_post_fkey" FOREIGN KEY ("associated_post") REFERENCES "LessonPost"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "LessonPost" ADD CONSTRAINT "LessonPost_associated_course_fkey" FOREIGN KEY ("associated_course") REFERENCES "Course"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "OverrideLessonPostVisibility" ADD CONSTRAINT "OverrideLessonPostVisibility_allowed_user_fkey" FOREIGN KEY ("allowed_user") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "OverrideLessonPostVisibility" ADD CONSTRAINT "OverrideLessonPostVisibility_allowed_post_id_fkey" FOREIGN KEY ("allowed_post_id") REFERENCES "LessonPost"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "OverrideLessonPostPermission" ADD CONSTRAINT "OverrideLessonPostPermission_allowed_editor_fkey" FOREIGN KEY ("allowed_editor") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "OverrideLessonPostPermission" ADD CONSTRAINT "OverrideLessonPostPermission_allowed_post_id_fkey" FOREIGN KEY ("allowed_post_id") REFERENCES "LessonPost"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "CourseInstructorMapping" ADD CONSTRAINT "CourseInstructorMapping_instructor_fkey" FOREIGN KEY ("instructor") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "CourseInstructorMapping" ADD CONSTRAINT "CourseInstructorMapping_course_id_fkey" FOREIGN KEY ("course_id") REFERENCES "Course"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "CourseLearnerMapping" ADD CONSTRAINT "CourseLearnerMapping_learner_fkey" FOREIGN KEY ("learner") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "CourseLearnerMapping" ADD CONSTRAINT "CourseLearnerMapping_course_id_fkey" FOREIGN KEY ("course_id") REFERENCES "Course"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "LiveClass" ADD CONSTRAINT "LiveClass_associated_course_fkey" FOREIGN KEY ("associated_course") REFERENCES "Course"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "OverrideLiveClassVisibility" ADD CONSTRAINT "OverrideLiveClassVisibility_learner_fkey" FOREIGN KEY ("learner") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "OverrideLiveClassVisibility" ADD CONSTRAINT "OverrideLiveClassVisibility_live_class_id_fkey" FOREIGN KEY ("live_class_id") REFERENCES "LiveClass"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "OverrideLiveClassAdminPermission" ADD CONSTRAINT "OverrideLiveClassAdminPermission_instructor_fkey" FOREIGN KEY ("instructor") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "OverrideLiveClassAdminPermission" ADD CONSTRAINT "OverrideLiveClassAdminPermission_live_class_id_fkey" FOREIGN KEY ("live_class_id") REFERENCES "LiveClass"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
